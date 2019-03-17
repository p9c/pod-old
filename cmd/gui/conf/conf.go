package conf

import (
	"fmt"

	"git.parallelcoin.io/dev/pod/cmd/gui/jdb"
)

type Conf struct {
	Interface InfConf    `json="interface"`
	Mining    MiningConf `json="mining"`
	Network   NetConf    `json="network"`
	Security  SecConf    `json="security"`
}

var VCF Conf = Conf{}

func (cf *Conf) ConfData() {

	if err := jdb.JDB.Read("conf", "interface", &cf.Interface); err != nil {

		fmt.Println("Error", err)
	}
	fmt.Println("Errosssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssr", cf.Interface)

	if err := jdb.JDB.Read("conf", "mining", &VCF.Mining); err != nil {

		fmt.Println("Error", err)
	}

	if err := jdb.JDB.Read("conf", "network", &VCF.Network); err != nil {

		fmt.Println("Error", err)
	}

	if err := jdb.JDB.Read("conf", "security", &VCF.Security); err != nil {

		fmt.Println("Error", err)
	}
}

func (cf *Conf) SaveInterfaceConf(lang string) {

	ICF := InfConf{

		Lang: lang,
	}
	jdb.JDB.Write("conf", "interface", ICF)
	fmt.Println("333333333sssssssssssssssssssssssssssssssssssssssssr", ICF)
	fmt.Println("langlanglanglanglanglanglanglanglang", lang)

}
