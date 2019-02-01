var VRFSendForm = {
template: vuedata.data.pages.vrfsend,
props:{
  address:String,
  label:String,
  amount:String,
  timeout:220,
  },
  // components: {
  //   VRFSendForm
  // },
  methods: {
    sendForm() {
      rpcinterface.walletPassphrase(this.passphrase, this.timeout);
      rpcinterface.sendToAddress(this.address, this.label, this.amount);
     },
    },
  data() {
      return {
            passphrase:"",
            amount:0,
            timeout:220,
            isComponentModalActive: false,
      }
  }
}