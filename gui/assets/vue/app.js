Vue.use(Buefy);
Vue.use(EasyBar);
var app = new Vue({
  el: '#app',
  data () {
    return {
    vdt: vuedata,
    bcd: blockchaindata,
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
  // this.timer = setInterval(this.ref, 5000)
},
methods: {
  swapComponent: function(component){this.component = component;},
  ref: function() { blockchaindata.getBlockChainData(); },
  cancelAutoUpdate: function() { clearInterval(this.timer) }
},
beforeDestroy() {
clearInterval(this.timer)
}
});