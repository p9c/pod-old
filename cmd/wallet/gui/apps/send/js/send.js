var SendC = {
  template: vuedata.data.pages.send,
  props:{
     vicons:Object,
    vlng:Object,
  },
  
  data () {
    return {
      address:"",
      label:"",
      amount:"",
      timeout:220,
      transactions: blockchaindata.data.listallsendtransactions,
      isPaginated: true,
      isPaginationSimple: false,
      defaultSortDirection: 'asc',
      currentPage: 1,
      perPage: 5,
    }
  },
  components: {
},
  methods: {
    vrfSend() {
      this.$modal.open({
          parent: this,
          component: VRFSendForm,
          hasModalCard: true,
          props:{
            address:this.address,
            label:this.label,
            amount:this.amount,
            timeout:220,
            },
            onSubmit: () => this.$toast.open('Transaction sent!')
      })
  }
    },

}