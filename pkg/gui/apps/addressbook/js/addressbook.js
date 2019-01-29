var AddressBookC = {
  template: vuedata.data.pages.addressbook,
  props:{
    vicons:Object,
    vlng:Object,
  },
  data () {
    return {
      address:"",
      label:"",
      desc:"",
      ab: addressbook.data,
      isPaginated: true,
      isPaginationSimple: false,
      defaultSortDirection: 'asc',
      currentPage: 1,
      perPage: 10,
    }
  },
  components: {
},
created: function() {
  this.abook();
  this.timer = setInterval(this.abook, 500)
},
  methods: {
    labelForm: function() {
      addressbooklabel.addressBookLabelWrite(this.label, this.address, this.desc);
      },
    delForm: function() {
        addressbooklabel.AddressBookLabelDelete();
    },
    confirmCustomDelete() {
        this.$dialog.confirm({
          title: 'Deleting label',
          message: 'Are you sure you want to <b>delete</b> your label? This action cannot be undone.',
          confirmText: 'Delete Label',
          type: 'is-danger',
          hasIcon: true,
          // function() { addressbooklabel.AddressBookLabelDelete();
          onConfirm: () => this.$toast.open('Label deleted!')
      })
    },
    abook: function() { addressbook.addressBookData(); },
    cancelAutoUpdate: function() { clearInterval(this.timer) }
  },
    beforeDestroy() {
    clearInterval(this.timer)
  }
}