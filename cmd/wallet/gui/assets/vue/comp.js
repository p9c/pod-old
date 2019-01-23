const Title = {
  template:`<div id="title"><p>{{title}}</p></div>`,
  props: ['title'],
};

var HomeC = {
  template: vuedata.data.pages.home,
  props:{
    vpages:Object,
    vicons:Object,
  }
}

var SendC = {
  template: vuedata.data.pages.send,
  props:{
    vpages:Object,
    vicons:Object,
  }
}


var ReceiveC = {
  template: vuedata.data.pages.receive,
  props:{
    vpages:Object,
    vicons:Object,
  }
}

var AddressBookC = {
  template: vuedata.data.pages.addressbook,
  props:{
    vpages:Object,
    vicons:Object,
  }
}

var HistoryC = {
  template: vuedata.data.pages.history,
  props:{
    vpages:Object,
    vicons:Object,
  }
}
