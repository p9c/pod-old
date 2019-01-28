var ReceiveC = {
  template: vuedata.data.pages.receive,
  props:{
    vicons:Object,
    vlng:Object,
  },
  data () {
    return {
      address:"",
      label:"",
      desc:"",
      amount:"",
      reqpay: reqpays.data,
      isPaginated: true,
      isPaginationSimple: false,
      defaultSortDirection: 'asc',
      currentPage: 1,
      perPage: 10,
    }
  },
  components: {
},
  methods: {
    reqForm: function() {
      reqpay.requestedPaymentWrite(this.label, this.address, this.amount, this.desc);
     },
  confirmCustomDelete() {
      this.$dialog.confirm({
          title: 'Deleting label',
          message: 'Are you sure you want to <b>delete</b> your label? This action cannot be undone.',
          confirmText: 'Delete Label',
          type: 'is-danger',
          hasIcon: true,
          onConfirm: () => this.$toast.open('Label deleted!')
      })
  }
}
}