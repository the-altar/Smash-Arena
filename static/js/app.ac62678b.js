(function(e){function t(t){for(var r,a,c=t[0],u=t[1],i=t[2],l=0,f=[];l<c.length;l++)a=c[l],Object.prototype.hasOwnProperty.call(o,a)&&o[a]&&f.push(o[a][0]),o[a]=0;for(r in u)Object.prototype.hasOwnProperty.call(u,r)&&(e[r]=u[r]);d&&d(t);while(f.length)f.shift()();return s.push.apply(s,i||[]),n()}function n(){for(var e,t=0;t<s.length;t++){for(var n=s[t],r=!0,a=1;a<n.length;a++){var c=n[a];0!==o[c]&&(r=!1)}r&&(s.splice(t--,1),e=u(u.s=n[0]))}return e}var r={},a={app:0},o={app:0},s=[];function c(e){return u.p+"js/"+({}[e]||e)+"."+{"chunk-6029506a":"cf704f6b","chunk-efac7da8":"a6e3608c"}[e]+".js"}function u(t){if(r[t])return r[t].exports;var n=r[t]={i:t,l:!1,exports:{}};return e[t].call(n.exports,n,n.exports,u),n.l=!0,n.exports}u.e=function(e){var t=[],n={"chunk-6029506a":1,"chunk-efac7da8":1};a[e]?t.push(a[e]):0!==a[e]&&n[e]&&t.push(a[e]=new Promise(function(t,n){for(var r="css/"+({}[e]||e)+"."+{"chunk-6029506a":"00667f41","chunk-efac7da8":"10012067"}[e]+".css",o=u.p+r,s=document.getElementsByTagName("link"),c=0;c<s.length;c++){var i=s[c],l=i.getAttribute("data-href")||i.getAttribute("href");if("stylesheet"===i.rel&&(l===r||l===o))return t()}var f=document.getElementsByTagName("style");for(c=0;c<f.length;c++){i=f[c],l=i.getAttribute("data-href");if(l===r||l===o)return t()}var d=document.createElement("link");d.rel="stylesheet",d.type="text/css",d.onload=t,d.onerror=function(t){var r=t&&t.target&&t.target.src||o,s=new Error("Loading CSS chunk "+e+" failed.\n("+r+")");s.code="CSS_CHUNK_LOAD_FAILED",s.request=r,delete a[e],d.parentNode.removeChild(d),n(s)},d.href=o;var m=document.getElementsByTagName("head")[0];m.appendChild(d)}).then(function(){a[e]=0}));var r=o[e];if(0!==r)if(r)t.push(r[2]);else{var s=new Promise(function(t,n){r=o[e]=[t,n]});t.push(r[2]=s);var i,l=document.createElement("script");l.charset="utf-8",l.timeout=120,u.nc&&l.setAttribute("nonce",u.nc),l.src=c(e);var f=new Error;i=function(t){l.onerror=l.onload=null,clearTimeout(d);var n=o[e];if(0!==n){if(n){var r=t&&("load"===t.type?"missing":t.type),a=t&&t.target&&t.target.src;f.message="Loading chunk "+e+" failed.\n("+r+": "+a+")",f.name="ChunkLoadError",f.type=r,f.request=a,n[1](f)}o[e]=void 0}};var d=setTimeout(function(){i({type:"timeout",target:l})},12e4);l.onerror=l.onload=i,document.head.appendChild(l)}return Promise.all(t)},u.m=e,u.c=r,u.d=function(e,t,n){u.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:n})},u.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},u.t=function(e,t){if(1&t&&(e=u(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var n=Object.create(null);if(u.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var r in e)u.d(n,r,function(t){return e[t]}.bind(null,r));return n},u.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return u.d(t,"a",t),t},u.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},u.p="/",u.oe=function(e){throw console.error(e),e};var i=window["webpackJsonp"]=window["webpackJsonp"]||[],l=i.push.bind(i);i.push=t,i=i.slice();for(var f=0;f<i.length;f++)t(i[f]);var d=l;s.push([0,"chunk-vendors"]),n()})({0:function(e,t,n){e.exports=n("56d7")},4483:function(e,t,n){},"56d7":function(e,t,n){"use strict";n.r(t);n("cadf"),n("551c"),n("f751"),n("097d");var r=n("2b0e"),a=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{attrs:{id:"app"}},[n("Nav"),n("router-view",{staticClass:"router-window"})],1)},o=[],s=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{attrs:{id:"nav"}},[n("TeamView"),n("input",{directives:[{name:"model",rawName:"v-model",value:e.username,expression:"username"}],attrs:{type:"text"},domProps:{value:e.username},on:{input:function(t){t.target.composing||(e.username=t.target.value)}}}),n("button",{on:{click:e.start_game}},[e._v("Start game")])],1)},c=[],u=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("section",[n("div",{staticClass:"team_view"},e._l(e.selected,function(t,r){return n("img",{key:r,staticClass:"icon icon--small icon--round team_view__char",attrs:{src:"https://smasharena.herokuapp.com/character/profile/"+t.ID+".jpg",alt:t.Name},on:{dblclick:function(t){return e.remove_char(r)}}})}),0)])},i=[],l={data:function(){return{username:""}},computed:{selected:function(){return this.$store.getters["character/selected"]}},watch:{username:function(e){this.$store.dispatch("user/setUser",e),this.username=e}},methods:{remove_char:function(e){this.$store.dispatch("character/remove_char",e)}}},f=l,d=(n("5ae5"),n("2877")),m=Object(d["a"])(f,u,i,!1,null,"7fc1864e",null),h=m.exports,p={data:function(){return{username:"",ws:null}},components:{TeamView:h},methods:{start_game:function(){var e=this;if(this.username.length){var t=this.$store.getters["character/selectedID"],n=this;3===t.length&&this.$http.post("/newgame",{userID:this.username,teamID:t}).then(function(t){t.data&&e.$store.dispatch("socket/create",{url:"wss://smasharena.herokuapp.com/arena/"+e.username,Vue:n})})}}}},g=p,v=Object(d["a"])(g,s,c,!1,null,"2436b2c3",null),b=v.exports,w={components:{Nav:b}},k=w,y=Object(d["a"])(k,a,o,!1,null,"752665d4",null),_=y.exports,S=n("bc3a"),O=n.n(S),j=n("8c4f");r["a"].use(j["a"]);var $=new j["a"]({mode:"history",base:"/",routes:[{path:"/",name:"home",component:function(){return n.e("chunk-efac7da8").then(n.bind(null,"1c62"))}},{path:"/arena",name:"arena",component:function(){return n.e("chunk-6029506a").then(n.bind(null,"2761"))}}]}),x=n("2f62"),E=(n("7f7f"),n("6b54"),{namespaced:!0,state:{roster:[],selected:[]},getters:{roster:function(e){return e.roster},selected:function(e){return e.selected},selectedID:function(e){return e.selected.map(function(e){return e.ID.toString()})}},mutations:{create:function(e,t){e.roster=t},select:function(e,t){var n=e.selected.length;n<3&&e.selected.push(e.roster.slice(t,t+1)[0])},remove:function(e,t){console.log(e.selected[t].name),console.log(e.selected[t].profile)}},actions:{fetch_all:function(e){O.a.get("/character").then(function(t){e.commit("create",t.data.roster)})},select_char:function(e,t){e.commit("select",t)},remove_char:function(e,t){e.commit("remove",t)}}}),P={namespaced:!0,state:{user:""},getters:{user:function(e){return e.user}},mutations:{setUser:function(e,t){e.user=t}},actions:{setUser:function(e,t){e.commit("setUser",t)}}},C={namespaced:!0,state:{gameState:{},foes:{char1:{id:1,health:100,skills:{}},char2:{id:2,health:100,skills:{}},char3:{id:3,health:100,skills:{}}},friends:{char1:{id:1,health:100,skills:{skill1:{id:1}}},char2:{id:2,health:50,skills:{skill1:{id:1}}},char3:{id:3,health:80,skills:{skill1:{id:1}}}},opponent:null},getters:{gameState:function(e){return e.gameState},friends:function(e){return e.friends},foes:function(e){return e.foes},opponent:function(e){return e.opponent}},mutations:{updateGameState:function(e,t){e.gameState=t},updateFriends:function(e,t){e.friends=t},updateFoes:function(e,t){e.foes=t}},actions:{update:function(e,t){console.log(t),e.commit("updateGameState",t.gameState),e.commit("updateFoes",t.gameState.foes),e.commit("updateFriends",t.gameState.friends)}}},N={namespaced:!0,state:{ws:null},getters:{ws:function(e){return e.ws}},mutations:{create:function(e,t){e.ws=new WebSocket(t.url),e.ws.onmessage=function(e){var n=JSON.parse(e.data);console.log(n),t.Vue.$store.dispatch("game/update",n),t.Vue.$router.push("arena")}}},actions:{create:function(e,t){console.log(t),e.commit("create",t)}}};r["a"].use(x["a"]);var T=new x["a"].Store({modules:{character:E,user:P,game:C,socket:N}});n("6418");r["a"].prototype.$http=O.a,r["a"].config.productionTip=!1,new r["a"]({router:$,store:T,render:function(e){return e(_)}}).$mount("#app")},"5ae5":function(e,t,n){"use strict";var r=n("4483"),a=n.n(r);a.a},6418:function(e,t,n){}});
//# sourceMappingURL=app.ac62678b.js.map