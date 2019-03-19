var VueVirtualScroller=function(e,t){"use strict";t=t&&t.hasOwnProperty("default")?t.default:t;var i={itemsLimit:1e3};var n=void 0;function r(){r.init||(r.init=!0,n=-1!==function(){var e=window.navigator.userAgent,t=e.indexOf("MSIE ");if(t>0)return parseInt(e.substring(t+5,e.indexOf(".",t)),10);if(e.indexOf("Trident/")>0){var i=e.indexOf("rv:");return parseInt(e.substring(i+3,e.indexOf(".",i)),10)}var n=e.indexOf("Edge/");return n>0?parseInt(e.substring(n+5,e.indexOf(".",n)),10):-1}())}var o={render:function(){var e=this.$createElement;return(this._self._c||e)("div",{staticClass:"resize-observer",attrs:{tabindex:"-1"}})},staticRenderFns:[],_scopeId:"data-v-b329ee4c",name:"resize-observer",methods:{compareAndNotify:function(){this._w===this.$el.offsetWidth&&this._h===this.$el.offsetHeight||(this._w=this.$el.offsetWidth,this._h=this.$el.offsetHeight,this.$emit("notify"))},addResizeHandlers:function(){this._resizeObject.contentDocument.defaultView.addEventListener("resize",this.compareAndNotify),this.compareAndNotify()},removeResizeHandlers:function(){this._resizeObject&&this._resizeObject.onload&&(!n&&this._resizeObject.contentDocument&&this._resizeObject.contentDocument.defaultView.removeEventListener("resize",this.compareAndNotify),delete this._resizeObject.onload)}},mounted:function(){var e=this;r(),this.$nextTick(function(){e._w=e.$el.offsetWidth,e._h=e.$el.offsetHeight});var t=document.createElement("object");this._resizeObject=t,t.setAttribute("aria-hidden","true"),t.setAttribute("tabindex",-1),t.onload=this.addResizeHandlers,t.type="text/html",n&&this.$el.appendChild(t),t.data="about:blank",n||this.$el.appendChild(t)},beforeDestroy:function(){this.removeResizeHandlers()}};var s={version:"0.4.5",install:function(e){e.component("resize-observer",o),e.component("ResizeObserver",o)}},l=null;"undefined"!=typeof window?l=window.Vue:"undefined"!=typeof global&&(l=global.Vue),l&&l.use(s);var a="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(e){return typeof e}:function(e){return e&&"function"==typeof Symbol&&e.constructor===Symbol&&e!==Symbol.prototype?"symbol":typeof e},c=(function(){function e(e){this.value=e}function t(t){var i,n;function r(i,n){try{var s=t[i](n),l=s.value;l instanceof e?Promise.resolve(l.value).then(function(e){r("next",e)},function(e){r("throw",e)}):o(s.done?"return":"normal",s.value)}catch(e){o("throw",e)}}function o(e,t){switch(e){case"return":i.resolve({value:t,done:!0});break;case"throw":i.reject(t);break;default:i.resolve({value:t,done:!1})}(i=i.next)?r(i.key,i.arg):n=null}this._invoke=function(e,t){return new Promise(function(o,s){var l={key:e,arg:t,resolve:o,reject:s,next:null};n?n=n.next=l:(i=n=l,r(e,t))})},"function"!=typeof t.return&&(this.return=void 0)}"function"==typeof Symbol&&Symbol.asyncIterator&&(t.prototype[Symbol.asyncIterator]=function(){return this}),t.prototype.next=function(e){return this._invoke("next",e)},t.prototype.throw=function(e){return this._invoke("throw",e)},t.prototype.return=function(e){return this._invoke("return",e)}}(),function(e,t){if(!(e instanceof t))throw new TypeError("Cannot call a class as a function")}),u=function(){function e(e,t){for(var i=0;i<t.length;i++){var n=t[i];n.enumerable=n.enumerable||!1,n.configurable=!0,"value"in n&&(n.writable=!0),Object.defineProperty(e,n.key,n)}}return function(t,i,n){return i&&e(t.prototype,i),n&&e(t,n),t}}(),h=function(e){if(Array.isArray(e)){for(var t=0,i=Array(e.length);t<e.length;t++)i[t]=e[t];return i}return Array.from(e)};var d=function(){function e(t,i,n){c(this,e),this.el=t,this.observer=null,this.frozen=!1,this.createObserver(i,n)}return u(e,[{key:"createObserver",value:function(e,t){var i,n,r,o,s,l,a,c=this;(this.observer&&this.destroyObserver(),this.frozen)||(this.options="function"==typeof(i=e)?{callback:i}:i,this.callback=this.options.callback,this.callback&&this.options.throttle&&(this.callback=(n=this.callback,r=this.options.throttle,o=void 0,s=void 0,l=void 0,(a=function(e){for(var t=arguments.length,i=Array(t>1?t-1:0),a=1;a<t;a++)i[a-1]=arguments[a];l=i,o&&e===s||(s=e,clearTimeout(o),o=setTimeout(function(){n.apply(void 0,[e].concat(h(l))),o=0},r))})._clear=function(){clearTimeout(o)},a)),this.oldResult=void 0,this.observer=new IntersectionObserver(function(e){var t=e[0];if(c.callback){var i=t.isIntersecting&&t.intersectionRatio>=c.threshold;if(i===c.oldResult)return;c.oldResult=i,c.callback(i,t),i&&c.options.once&&(c.frozen=!0,c.destroyObserver())}},this.options.intersection),t.context.$nextTick(function(){c.observer.observe(c.el)}))}},{key:"destroyObserver",value:function(){this.observer&&(this.observer.disconnect(),this.observer=null),this.callback&&this.callback._clear&&(this.callback._clear(),this.callback=null)}},{key:"threshold",get:function(){return this.options.intersection&&this.options.intersection.threshold||0}}]),e}();function f(e,t,i){var n=t.value;if("undefined"==typeof IntersectionObserver)console.warn("[vue-observe-visibility] IntersectionObserver API is not available in your browser. Please install this polyfill: https://github.com/w3c/IntersectionObserver/tree/master/polyfill");else{var r=new d(e,n,i);e._vue_visibilityState=r}}var v={bind:f,update:function(e,t,i){var n=t.value;if(!function e(t,i){if(t===i)return!0;if("object"===(void 0===t?"undefined":a(t))){for(var n in t)if(!e(t[n],i[n]))return!1;return!0}return!1}(n,t.oldValue)){var r=e._vue_visibilityState;r?r.createObserver(n,i):f(e,{value:n},i)}},unbind:function(e){var t=e._vue_visibilityState;t&&(t.destroyObserver(),delete e._vue_visibilityState)}};var p={version:"0.4.3",install:function(e){e.directive("observe-visibility",v)}},m=null;"undefined"!=typeof window?m=window.Vue:"undefined"!=typeof global&&(m=global.Vue),m&&m.use(p);var y="undefined"!=typeof window?window:"undefined"!=typeof global?global:"undefined"!=typeof self?self:{};var g,b=(function(e){var t,i;t=y,i=function(){var e=/(auto|scroll)/,t=function(e,i){return null===e.parentNode?i:t(e.parentNode,i.concat([e]))},i=function(e,t){return getComputedStyle(e,null).getPropertyValue(t)},n=function(t){return e.test(function(e){return i(e,"overflow")+i(e,"overflow-y")+i(e,"overflow-x")}(t))};return function(e){if(e instanceof HTMLElement||e instanceof SVGElement){for(var i=t(e.parentNode,[]),r=0;r<i.length;r+=1)if(n(i[r]))return i[r];return document.scrollingElement||document.documentElement}}},e.exports?e.exports=i():t.Scrollparent=i()}(g={exports:{}},g.exports),g.exports),_=!1;if("undefined"!=typeof window){_=!1;try{var w=Object.defineProperty({},"passive",{get:function(){_=!0}});window.addEventListener("test",null,w)}catch(e){}}var $=0,S={render:function(){var e=this,t=e.$createElement,i=e._self._c||t;return i("div",{directives:[{name:"observe-visibility",rawName:"v-observe-visibility",value:e.handleVisibilityChange,expression:"handleVisibilityChange"}],staticClass:"vue-recycle-scroller",class:{ready:e.ready,"page-mode":e.pageMode},on:{"&scroll":function(t){return e.handleScroll(t)}}},[e._t("before-container"),e._v(" "),i("div",{ref:"wrapper",staticClass:"vue-recycle-scroller__item-wrapper",style:{height:e.totalHeight+"px"}},e._l(e.pool,function(t){return i("div",{key:t.nr.id,staticClass:"vue-recycle-scroller__item-view",class:{hover:e.hoverKey===t.nr.key},style:e.ready?{transform:"translateY("+t.top+"px)"}:null,on:{mouseenter:function(i){e.hoverKey=t.nr.key},mouseleave:function(t){e.hoverKey=null}}},[e._t("default",null,{item:t.item,index:t.nr.index,active:t.nr.used})],2)}),0),e._v(" "),e._t("after-container"),e._v(" "),i("ResizeObserver",{on:{notify:e.handleResize}})],2)},staticRenderFns:[],name:"RecycleScroller",mixins:[{components:{ResizeObserver:o},directives:{ObserveVisibility:v},props:{items:{type:Array,required:!0},itemHeight:{type:Number,default:null},minItemHeight:{type:[Number,String],default:null},heightField:{type:String,default:"height"},typeField:{type:String,default:"type"},keyField:{type:String,default:"id"},buffer:{type:Number,default:200},pageMode:{type:Boolean,default:!1},prerender:{type:Number,default:0},emitUpdate:{type:Boolean,default:!1}},computed:{heights:function(){if(null===this.itemHeight){for(var e={"-1":{accumulator:0}},t=this.items,i=this.heightField,n=this.minItemHeight,r=0,o=void 0,s=0,l=t.length;s<l;s++)r+=o=t[s][i]||n,e[s]={accumulator:r,height:o};return e}return[]}},beforeDestroy:function(){this.removeListeners()},methods:{getListenerTarget:function(){var e=b(this.$el);return e!==window.document.documentElement&&e!==window.document.body||(e=window),e},getScroll:function(){var e=this.$el,t=void 0;if(this.pageMode){var i=e.getBoundingClientRect(),n=-i.top,r=window.innerHeight;n<0&&(r+=n,n=0),n+r>i.height&&(r=i.height-n),t={top:n,bottom:n+r}}else t={top:e.scrollTop,bottom:e.scrollTop+e.clientHeight};return t},applyPageMode:function(){this.pageMode?this.addListeners():this.removeListeners()},addListeners:function(){this.listenerTarget=this.getListenerTarget(),this.listenerTarget.addEventListener("scroll",this.handleScroll,!!_&&{passive:!0}),this.listenerTarget.addEventListener("resize",this.handleResize)},removeListeners:function(){this.listenerTarget&&(this.listenerTarget.removeEventListener("scroll",this.handleScroll),this.listenerTarget.removeEventListener("resize",this.handleResize),this.listenerTarget=null)},scrollToItem:function(e){var t=void 0;t=null===this.itemHeight?e>0?this.heights[e-1].accumulator:0:e*this.itemHeight,this.scrollToPosition(t)},scrollToPosition:function(e){this.$el.scrollTop=e},itemsLimitError:function(){var e=this;throw setTimeout(function(){console.log("It seems the scroller element isn't scrolling, so it tries to render all the items at once.","Scroller:",e.$el),console.log("Make sure the scroller has a fixed height and 'overflow-y' set to 'auto' so it can scroll correctly and only render the items visible in the scroll viewport.")}),new Error("Rendered items limit reached")}}}],data:function(){return{pool:[],totalHeight:0,ready:!1,hoverKey:null}},watch:{items:function(){this.updateVisibleItems(!0)},pageMode:function(){this.applyPageMode(),this.updateVisibleItems(!1)},heights:{handler:function(){this.updateVisibleItems(!1)},deep:!0}},created:function(){this.$_startIndex=0,this.$_endIndex=0,this.$_views=new Map,this.$_unusedViews=new Map,this.$_scrollDirty=!1,this.$isServer&&this.updateVisibleItems(!1)},mounted:function(){var e=this;this.applyPageMode(),this.$nextTick(function(){e.updateVisibleItems(!0),e.ready=!0})},methods:{addView:function(e,t,i,n,r){var o={item:i,top:0},s={id:$++,index:t,used:!0,key:n,type:r};return Object.defineProperty(o,"nr",{configurable:!1,value:s}),e.push(o),o},unuseView:function(e){var t=arguments.length>1&&void 0!==arguments[1]&&arguments[1],i=this.$_unusedViews,n=e.nr.type,r=i.get(n);r||(r=[],i.set(n,r)),r.push(e),t||(e.nr.used=!1,e.top=-9999,this.$_views.delete(e.nr.key))},handleResize:function(){this.$emit("resize"),this.ready&&this.updateVisibleItems(!1)},handleScroll:function(e){var t=this;this.$_scrollDirty||(this.$_scrollDirty=!0,requestAnimationFrame(function(){t.$_scrollDirty=!1,t.updateVisibleItems(!1).continuous||(clearTimeout(t.$_refreshTimout),t.$_refreshTimout=setTimeout(t.handleScroll,100))}))},handleVisibilityChange:function(e,t){var i=this;this.ready&&(e||0!==t.boundingClientRect.width||0!==t.boundingClientRect.height?(this.$emit("visible"),requestAnimationFrame(function(){i.updateVisibleItems(!1)})):this.$emit("hidden"))},updateVisibleItems:function(e){var t=this.itemHeight,n=this.typeField,r=this.keyField,o=this.items,s=o.length,l=this.heights,a=this.$_views,c=this.$_unusedViews,u=this.pool,h=void 0,d=void 0,f=void 0;if(s)if(this.$isServer)h=0,d=this.prerender,f=null;else{var v=this.getScroll(),p=this.buffer;if(v.top-=p,v.bottom+=p,null===t){var m=0,y=s-1,g=~~(s/2),b=void 0;do{b=g,l[g].accumulator<v.top?m=g:g<s-1&&l[g+1].accumulator>v.top&&(y=g),g=~~((m+y)/2)}while(g!==b);for(g<0&&(g=0),h=g,f=l[s-1].accumulator,d=g;d<s&&l[d].accumulator<v.bottom;d++);-1===d?d=o.length-1:++d>s&&(d=s)}else h=~~(v.top/t),d=Math.ceil(v.bottom/t),h<0&&(h=0),d>s&&(d=s),f=s*t}else h=d=f=0;d-h>i.itemsLimit&&this.itemsLimitError(),this.totalHeight=f;var _=void 0,w=h<=this.$_endIndex&&d>=this.$_startIndex,$=void 0;if(this.$_continuous!==w){if(w){a.clear(),c.clear();for(var S=0,x=u.length;S<x;S++)_=u[S],this.unuseView(_)}this.$_continuous=w}else if(w)for(var z=0,V=u.length;z<V;z++)(_=u[z]).nr.used&&(e&&(_.nr.index=o.findIndex(function(e){return r?e[r]===_.item[r]:e===_.item})),(-1===_.nr.index||_.nr.index<h||_.nr.index>=d)&&this.unuseView(_));w||($=new Map);for(var I=void 0,D=void 0,k=void 0,O=void 0,R=h;R<d;R++){I=o[R];var T=r?I[r]:I;_=a.get(T),t||l[R].height?(_?(_.nr.used=!0,_.item=I):(D=I[n],w?(k=c.get(D))&&k.length?((_=k.pop()).item=I,_.nr.used=!0,_.nr.index=R,_.nr.key=T,_.nr.type=D):_=this.addView(u,R,I,T,D):(k=c.get(D),O=$.get(D)||0,k&&O<k.length?((_=k[O]).item=I,_.nr.used=!0,_.nr.index=R,_.nr.key=T,_.nr.type=D,$.set(D,O+1)):(_=this.addView(u,R,I,T,D),this.unuseView(_,!0)),O++),a.set(T,_)),_.top=null===t?l[R-1].accumulator:R*t):_&&this.unuseView(_)}return this.$_startIndex=h,this.$_endIndex=d,this.emitUpdate&&this.$emit("update",h,d),{continuous:w}}}},x={render:function(){var e=this,t=e.$createElement,i=e._self._c||t;return i("RecycleScroller",e._g(e._b({ref:"scroller",attrs:{items:e.itemsWithHeight,"min-item-height":e.minItemHeight},on:{resize:e.onScrollerResize,visible:e.onScrollerVisible},scopedSlots:e._u([{key:"default",fn:function(t){var i=t.item,n=t.index,r=t.active;return[e._t("default",null,null,{item:i.item,index:n,active:r,itemWithHeight:i})]}}])},"RecycleScroller",e.$attrs,!1),e.listeners),[i("template",{slot:"before-container"},[e._t("before-container")],2),e._v(" "),i("template",{slot:"after-container"},[e._t("after-container")],2)],2)},staticRenderFns:[],name:"DynamicScroller",components:{RecycleScroller:S},inheritAttrs:!1,provide:function(){return{vscrollData:this.vscrollData,vscrollBus:this}},props:{items:{type:Array,required:!0},minItemHeight:{type:[Number,String],required:!0},keyField:{type:String,default:"id"}},data:function(){return{vscrollData:{active:!0,heights:{},keyField:this.keyField}}},computed:{itemsWithHeight:function(){for(var e=[],t=this.items,i=this.keyField,n=this.vscrollData.heights,r=0;r<t.length;r++){var o=t[r],s=o[i];e.push({item:o,id:s,height:n[s]||0})}return e},listeners:function(){var e={};for(var t in this.$listeners)"resize"!==t&&"visible"!==t&&(e[t]=this.$listeners[t]);return e}},watch:{items:"forceUpdate"},created:function(){this.$_updates=[]},mounted:function(){var e=this.$refs.scroller,t=this.getSize(e);this._scrollerWidth=t.width},activated:function(){this.vscrollData.active=!0},deactivated:function(){this.vscrollData.active=!1},methods:{onScrollerResize:function(){this.$refs.scroller&&this.forceUpdate(),this.$emit("resize")},onScrollerVisible:function(){this.$emit("vscroll:update",{force:!1}),this.$emit("visible")},forceUpdate:function(){this.vscrollData.heights={},this.$emit("vscroll:update",{force:!0})},getSize:function(e){return e.$el.getBoundingClientRect()},scrollToItem:function(e){var t=this.$refs.scroller;t&&t.scrollToItem(e)}}},z={name:"DynamicScrollerItem",inject:["vscrollData","vscrollBus"],props:{item:{type:Object,required:!0},watchData:{type:Boolean,default:!1},active:{type:Boolean,required:!0},sizeDependencies:{type:[Array,Object],default:null},emitResize:{type:Boolean,default:!1},tag:{type:String,default:"div"}},computed:{id:function(){return this.item[this.vscrollData.keyField]},height:function(){return this.vscrollData.heights[this.id]||0}},watch:{watchData:"updateWatchData",id:function(){this.height||this.onDataUpdate()},active:function(e){e&&this.$_pendingVScrollUpdate&&this.updateSize()}},created:function(){var e=this;if(!this.$isServer){this.$_forceNextVScrollUpdate=!1,this.updateWatchData();var t=function(t){e.$watch(function(){return e.sizeDependencies[t]},e.onDataUpdate)};for(var i in this.sizeDependencies)t(i);this.vscrollBus.$on("vscroll:update",this.onVscrollUpdate),this.vscrollBus.$on("vscroll:update-size",this.onVscrollUpdateSize)}},mounted:function(){this.vscrollData.active&&this.updateSize()},beforeDestroy:function(){this.vscrollBus.$off("vscroll:update",this.onVscrollUpdate),this.vscrollBus.$off("vscroll:update-size",this.onVscrollUpdateSize)},methods:{updateSize:function(){this.active&&this.vscrollData.active?this.$_pendingSizeUpdate||(this.$_pendingSizeUpdate=!0,this.$_forceNextVScrollUpdate=!1,this.$_pendingVScrollUpdate=!1,this.active&&this.vscrollData.active&&this.computeSize(this.id)):this.$_forceNextVScrollUpdate=!0},getSize:function(){return this.$el.getBoundingClientRect()},updateWatchData:function(){var e=this;this.watchData?this.$_watchData=this.$watch("data",function(){e.onDataUpdate()},{deep:!0}):this.$_watchData&&(this.$_watchData(),this.$_watchData=null)},onVscrollUpdate:function(e){var t=e.force;!this.active&&t&&(this.$_pendingVScrollUpdate=!0),(this.$_forceNextVScrollUpdate||t||!this.height)&&this.updateSize()},onDataUpdate:function(){this.updateSize()},computeSize:function(e){var t=this;this.$nextTick(function(){if(t.id===e){var i=t.getSize();i.height&&t.height!==i.height&&(t.$set(t.vscrollData.heights,t.id,i.height),t.emitResize&&t.$emit("resize",t.id))}t.$_pendingSizeUpdate=!1})}},render:function(e){return e(this.tag,this.$slots.default)}};var V={version:"1.0.0-beta.4",install:function(e,t){var n=Object.assign({},{installComponents:!0,componentsPrefix:""},t);for(var r in n)void 0!==n[r]&&(i[r]=n[r]);n.installComponents&&function(e,t){e.component(t+"recycle-scroller",S),e.component(t+"RecycleScroller",S),e.component(t+"dynamic-scroller",x),e.component(t+"DynamicScroller",x),e.component(t+"dynamic-scroller-item",z),e.component(t+"DynamicScrollerItem",z)}(e,n.componentsPrefix)}},I=null;return"undefined"!=typeof window?I=window.Vue:"undefined"!=typeof global&&(I=global.Vue),I&&I.use(V),e.RecycleScroller=S,e.DynamicScroller=x,e.DynamicScrollerItem=z,e.default=V,e.IdState=function(){var e=(arguments.length>0&&void 0!==arguments[0]?arguments[0]:{}).idProp,i=void 0===e?function(e){return e.item.id}:e,n={},r=new t({data:function(){return{store:n}}});return{data:function(){return{idState:null}},created:function(){var e=this;this.$_id=null,this.$_getId="function"==typeof i?function(){return i.call(e,e)}:function(){return e[i]},this.$watch(this.$_getId,{handler:function(e){var t=this;this.$nextTick(function(){t.$_id=e})},immediate:!0}),this.$_updateIdState()},beforeUpdate:function(){this.$_updateIdState()},methods:{$_idStateInit:function(e){var t=this.$options.idState;if("function"==typeof t){var i=t.call(this,this);return r.$set(n,e,i),this.$_id=e,i}throw new Error("[mixin IdState] Missing `idState` function on component definition.")},$_updateIdState:function(){var e=this.$_getId();null==e&&console.warn("No id found for IdState with idProp: '"+i+"'."),e!==this.$_id&&(n[e]||this.$_idStateInit(e),this.idState=n[e])}}}},e}({},Vue);