package nknovh_engine

import (
	"io/ioutil"
	"strconv"
	"os"
	nkn "github.com/nknorg/nkn-sdk-go"
)

type Nknsdk struct {
	Wallet *nkn.Wallet
}

func (o *NKNOVH) walletCreate() error {
	if _, err := os.Stat("external/wallet.json"); err == nil {
		if _, err := os.Stat("external/wallet.pswd"); err == nil {
			return nil
		}
	}

	account, err := nkn.NewAccount(nil)
	if err != nil {
		return err
	}

	wpswd := RandBytes(32)

	wallet, err := nkn.NewWallet(account, &nkn.WalletConfig{Password: string(wpswd)})
	if err != nil {
		return err
	}
	walletJSON, err := wallet.ToJSON()
	if err != nil {
		return err
	}

	wFile, err := os.OpenFile("external/wallet.json", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer wFile.Close()
	pFile, err := os.OpenFile("external/wallet.pswd", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer pFile.Close()

	if _, err := wFile.WriteString(walletJSON); err != nil {
		return err
	}
	if _, err := pFile.Write(wpswd); err != nil {
		return err
	}

	return nil
}

func (o *NKNOVH) nknConnect() error {
	data, err := ioutil.ReadFile("external/wallet.json")
	if err != nil {
		return err
	}
	wpswd, err := ioutil.ReadFile("external/wallet.pswd")
	if err != nil {
		return err
	}
	wallet, err := nkn.WalletFromJSON(string(data), &nkn.WalletConfig{Password: string(wpswd)})
	if err != nil {
		return err
	}
	if err := nkn.VerifyWalletAddress(wallet.Address()); err != nil {
		return err
	}
	if err := wallet.VerifyPassword(string(wpswd)); err != nil {
		return err
	}
	t := new(Nknsdk)
	t.Wallet = wallet
	o.Nknsdk = t
	return nil
}

func (o *NKNOVH) walletPoll() error {
	go o.getPrices()
	if err := o.fetchBalances(); err != nil {
		return err
	}
	return nil
}

func (o *NKNOVH) fetchBalances() error {

	var (
		id uint
		nkn_wallet string
		db_balance float64
	)
	var wallet *nkn.Wallet = o.Nknsdk.Wallet

	//fetch wallets from the database
	rows, err := o.sql.stmt["main"]["selectWallets"].Query()
	if err != nil { 
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &nkn_wallet, &db_balance); err != nil {
			return err
		}
		balance, err := wallet.BalanceByAddress(nkn_wallet)
		if err != nil {
			o.log.Syslog("Query balance fail: " + err.Error(), "wallets")
			continue
		}
		float_balance, err := strconv.ParseFloat(balance.String(), 64);
		if err != nil {
			o.log.Syslog("Cannot transform string to float64: " + err.Error(), "wallets")
			continue
		}
		if float_balance != db_balance {
			if _, err1 := o.sql.stmt["main"]["updateWalletBalanceById"].Exec(&float_balance, &id); err1 != nil {
				o.log.Syslog("Stmt updateWalletBalanceById has returned an error: ("+err1.Error()+")", "sql")
				continue
			}
		}
	}
	return nil
}
