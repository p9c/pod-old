var PeersC = {
  template: vuedata.data.pages.peers,
  props:{
    vlng:Object,
    vicons:Object,
  },
    data () {
      return { 
        zoom: 2,
        center: [0, 0],
        rotation: 0,
        geolocPosition: undefined,
      }
    },
}