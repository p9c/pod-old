Vue.use(Buefy);
Vue.use(VueVirtualScroller);
Vue.use(VueLayers);
Vue.use(VueTerminal);
 
var app = new Vue({
  el: '#app',
  data () {
    return {
    vdt: vuedata,
    bcd: blockchaindata,
    lng: language,
    ab: addressbook,
    rpc:rpcinterface,
    msg:"Welcome!",
    err:"Welcome!",
    cnf:conf,
    // rpc: rpchandlers,
    // vpage: vdt.data.pages.home,
    timer: '',
    // component: Home,
    component: HomeC,
    updateAvailable: false,
  }
},
components: {
  HomeC,
  // SendC,
},
created: function() {
  this.ref();
  this.lang();
  this.getBlockCount();
  this.getInfo();
  this.config();
  this.timer = setInterval(this.ref, 500)
  this.timer = setInterval(this.lang, 500)
  // this.timer = setInterval(this.getBlockCount, 1500)
  // this.timer = setInterval(this.getInfo, 1500)

  // this.timer = setInterval(this.danger, 500)
  // this.timer = setInterval(this.adrbk, 500)
},
watch:{
  'rpc.data.MSG': function(newVal, oldVal) {
    if (newVal != oldVal){
        // this.msg = newVal
        this.warning(newVal);
      }
  },
  'rpc.data.ERR': function(newVal, oldVal) {
    if (newVal != oldVal){
        // this.msg = newVal
        this.danger(newVal);
      }
  },
  'rpc.data.BlockCount': function(newVal, oldVal) {
    if (newVal != oldVal){
        // this.msg = newVal
        this.warning("New block: "+ newVal);
      }
  }
  // rpc :{
    // handler: function(val) {
      //   if (val.data.MSG != this.msg){
      //     this.msg = val.data.MSG
      //     this.warning(val.data.MSG);
      //   }
      //   if (val.data.ERR != this.err){
      //     this.err = val.data.ERR
      //     this.danger(val.data.ERR);
      //   }
      // },
      // deep: true
  // }
},
methods: {
  // processForm: function() {
  //   console.log({ name: this.name, email: this.email });
  //   alert('Processing');
  // },
  // rpc: function() { rpchandlers},
  warning(val) {
    this.$toast.open({
        duration: 5000,
        message: val,
        type: 'is-warning',
        position: 'is-top',
        actionText: 'close',
        queue: false,
        onAction: () => {
            this.$toast.open({
                message: 'Closed',
                queue: false
            })
        }
    })
  },
  danger(val) {
    this.$toast.open({
        duration: 5000,
        message: val,
        type: 'is-danger',
        position: 'is-top',
        actionText: 'close',
        queue: false,
        onAction: () => {
            this.$toast.open({
                message: 'Closed',
                queue: false
            })
        }
    })
  },
  swapComponent: function(component){this.component = component;},
  ref: function() { blockchaindata.getInfoData(); },
  getBlockCount: function() { rpcinterface.getBlockCount(); },
  getInfo: function() { rpcinterface.getInfo(); },
  config: function() { conf.confData(); },
  lang: function() { language.languageData(conf.data.Interface.lang); },
  cancelAutoUpdate: function() { clearInterval(this.timer) }
},
beforeDestroy() {
clearInterval(this.timer)
}
});

  