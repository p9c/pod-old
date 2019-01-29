var ConsoleC = {
  template: vuedata.data.pages.console,
props: {
  vlng:Object,
  vicons:Object,
  vrpc:Object,
  commands: {
    type: Object,
    default: {
      help: `$ help
  $addmultisigaddress
	$createmultisig    
	$dumpprivkey       
	$getaccount        
	$getaccountaddress 
	$getaddressesbyaccount
	$getbalance        
	$getbestblockhash  
	$getblockcount     
	$getinfo
	$getnewaddress     
	$getrawchangeaddress  
	$getreceivedbyaccount 
	$getreceivedbyaddress 
	$gettransaction    
	$help              
	$importprivkey     
	$keypoolrefill     
	$listaccounts      
	$listlockunspent   
	$listreceivedbyaccount
	$listreceivedbyaddress
	$listsinceblock
	$listtransactions  
	$listunspent       
	$lockunspent       
	$sendfrom
	$sendmany          
	$sendtoaddress     
	$settxfee          
	$signmessage       
	$signrawtransaction
	$validateaddress       
	$verifymessage          
	$walletlock            
	$walletpassphrase
	$walletpassphrasechange
  $backupwallet
	$dumpwallet
	$getwalletinfo
	$importwallet
	$listaddressgroupings
	$encryptwallet
	$move
	$setaccount
	$createnewaccount
	$getbestblock
	$getunconfirmedbalance
	$listaddresstransactions
	$listalltransactions
	$renameaccount
  $walletislocked
  `,
  clear: "exec clear",
  getbalance:" this.vrpc.Getinfo.balance",
  ls: "artus.txt skills.txt cat.txt \n",
      cat: "Specify the file you want to open... \n",
      "cat artus.txt": `\nMy name is Artus Vranken, and I am an avid software developer.
I'm always open to new challenges and aim to improve my skills each day.\n`,
      "cat cat.txt": `
　 ／l、     Meow :3
ﾞ（ﾟ､ ｡ ７
　 l、ﾞ ~ヽ
　 じしf_, )ノ
`,
      "cat skills.txt": `
Frontend:
  JS/html/css ★★★★
  Angular     ★★★
  Vue         ★★★

Backend:
  Java        ★★★★
  Spring      ★★
  Node.js     ★★★★

A bit of both:
  TypeScript  ★★★

A bit of everything else:
  Bash/Sh     ★★★★
  Docker      ★★★★
    -compose  ★★★
  Blockchain  ★★★
`
    }
  },
  user: {
    type: String,
    default: "duo@parallelcoin"
  }
},

/**
 * Data that's being tracked by this component.
 */
data: function() {
  return {
    directory: "/~",
    suffix: "$",
    history: new Array(),
    historyIndex: 0,
    input: "",
    output: new Array(),
    inputId: Math.floor(Math.random() * 1000)
  };
},

/**
 * Computed values.
 */
computed: {
  prefix: function() {
    return `${this.user}${this.directory} ${this.suffix}`;
  }
},

/**
 * Methods.
 */
methods: {
  /**
   * Check if the contenteditable span is in focus.
   */
  isFocused: function() {
    return document.activeElement.id == this.inputId;
  },

  /**
   * Set the focus to the contenteditable span.
   */
  focus: function() {
    while (document.activeElement.id != this.inputId) {
      document.getElementById(this.inputId).focus();
    }
  },

  /**
   * Perform specified actions when a keyUp event is fired.
   *
   * @param {eventArgs} e - The Event object.
   */
  keyUp: function(e) {
    switch (e.keyCode) {
      case 13:
        e.preventDefault();
        this.execute();
        break;

      case 38:
        e.preventDefault();
        this.previousHistory();
        break;
      case 40:
        e.preventDefault();
        this.nextHistory();
        break;
    }

    this.updateInputValue();
  },

  /**
   * Update the "input" data-field based on the contents of the terminal input field.
   */
  updateInputValue: function() {
    this.input = document.getElementById(this.inputId).innerHTML;
  },

  /**
   * Update the terminal input field based on the contents of the "input" data-field.
   */
  updateFieldValue: function() {
    document.getElementById(this.inputId).innerHTML = this.input;
  },

  /**
   * Execute functions entered by the user, based on the "commands" data-object.
   */
  execute: function() {
    let tempInput = this.input.replace("<br>", "");
    tempInput = tempInput.replace("<div>", "");
    tempInput = tempInput.replace("</div>", "");
    this.historyIndex = 0;
    this.history.unshift(tempInput);

    let tempOutput = this.commands[tempInput];

    if (typeof tempOutput == "undefined")
      tempOutput = `Couldn't find command: ${tempInput}`;

    switch (tempOutput) {
      case "exec clear":
        this.clear();
        return;
        break;
    }

    this.output.push(`${this.prefix} ${tempInput}`);
    this.output.push(tempOutput);

    document.getElementById(this.inputId).innerHTML = "";
    this.input = "";

    Vue.nextTick(function() {
      document.getElementById("vue-terminal").scrollBy(0, 10000);
      document.getElementsByClassName("vue-terminal-input")[0].focus();
    });
  },

  /**
   * Load previous command from history.
   */
  previousHistory: function() {
    if (this.historyIndex + 1 > this.history.length) return;
    this.input = this.history[this.historyIndex++];
    this.updateFieldValue();
  },

  /**
   * Load next command from history.
   */
  nextHistory: function() {
    if (this.historyIndex - 1 < 0) return;
    this.input = this.history[this.historyIndex--];
    this.updateFieldValue();
  },

  /**
   * Clear the input, both in the view as in the Vue instance.
   */
  clearInput: function() {
    document.getElementById(this.inputId).innerHTML = "";
    this.input = "";
  },

  /**
   * Clear the whole terminal screen.
   */
  clear: function() {
    this.output = new Array();
    this.clearInput();
  }
}
};