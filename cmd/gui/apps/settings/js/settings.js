var Ifc = {
  template: vuedata.data.pages.ifc,
  data() {
    return {
      // lang: "",
    }
},
  methods: {
    saveInterfaceConfForm() {
      conf.saveInterfaceConf(conf.data.Interface.lang);
     },
    },
}
var Network = {
  template: vuedata.data.pages.network,
}
var Security = {
  template: vuedata.data.pages.security,
}
var Mining = {
  template: vuedata.data.pages.mining,
}

var MiningAAA = {
  template: vuedata.data.pages.miningaaa,
}



var SettingsC = {
  template: vuedata.data.pages.settings,
  data() {
      return {
        component: Ifc,
      }
  },
  props:{
    vlng:Object,
    settings:Object,
  },
  methods: {
    swapComponent: function(component){this.component = component;},
}
}