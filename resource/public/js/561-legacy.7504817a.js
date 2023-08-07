"use strict";(self["webpackChunkarkose_vue2_example"]=self["webpackChunkarkose_vue2_example"]||[]).push([[561],{273:function(e,t,o){o.d(t,{Z:function(){return d}});var n=function(){var e=this,t=e._self._c;return"inline"===e.mode?t("div",{attrs:{id:e.selector?.slice(1)}}):e._e()},r=[],s={name:"Arkose",props:{publicKey:{type:String,default:""},mode:{type:String,default:""},selector:{type:String,default:null},nonce:{type:String,default:""}},data(){return{scriptId:""}},methods:{removeScript(){const e=document.getElementById(this.scriptId);e&&e.remove()},loadScript(e,t){this.removeScript();const o=document.createElement("script");return o.id=this.scriptId,o.type="text/javascript",o.src=`https://client-api.arkoselabs.com/v2/${e}/api.js`,o.setAttribute("data-callback","setupEnforcement"),o.async=!0,o.defer=!0,t&&o.setAttribute("data-nonce",t),document.body.appendChild(o),o},setupEnforcement(e){window.myEnforcement=e,window.myEnforcement.setConfig({selector:this.selector,mode:this.mode,onReady:()=>{this.$emit("onReady"),window.myEnforcement.run()},onShown:()=>{this.$emit("onShown")},onShow:()=>{this.$emit("onShow")},onSuppress:()=>{this.$emit("onSuppress")},onCompleted:e=>{this.$emit("onCompleted",e.token)},onReset:()=>{this.$emit("onReset")},onHide:()=>{this.$emit("onHide")},onError:e=>{this.$emit("onError",e)},onFailed:e=>{this.$emit("onFailed",e)}})}},mounted(){this.scriptId=`arkose-script-${this.publicKey}`;const e=this.loadScript(this.publicKey,this.nonce);e.onload=()=>{console.log("Arkose API Script loaded"),window.setupEnforcement=this.setupEnforcement.bind(this)},e.onerror=()=>{console.log("Could not load the Arkose API Script!")}},destroyed(){window.myEnforcement&&delete window.myEnforcement,window.setupEnforcement&&delete window.setupEnforcement,this.removeScript()}},i=s,l=o(1),c=(0,l.Z)(i,n,r,!1,null,null,null),d=c.exports},561:function(e,t,o){o.r(t),o.d(t,{default:function(){return a}});var n=function(){var e=this,t=e._self._c;return t("div",[t("h2",[e._v("正在获取token2")]),t("Arkose",{attrs:{"public-key":e.publicKey,mode:"lightbox"},on:{onCompleted:function(t){return e.onCompleted(t)},onError:function(t){return e.onError(t)}}}),t("div",{attrs:{id:"token"}},[e._v(" "+e._s(e.arkoseToken)+" ")]),t("input",{attrs:{type:"submit",value:"Submit"},on:{click:function(t){return e.onSubmit()}}})],1)},r=[],s=o(887),i=o(273),l={name:"DashboardComponent",components:{Arkose:i.Z},data(){return{publicKey:"35536E1E-65B4-4D96-9D97-6ADB7EFF8147",arkoseToken:null}},methods:{mounted(){console.log("mounted")},onCompleted(e){this.arkoseToken=e,console.log("token",e),this.$http.post("/pushtoken",{token:e}).then((e=>{console.log(e);const t=this.$route.query.delay||5;console.log("delay",t),setTimeout((()=>{s.Z.replace({path:"/",query:{delay:t}})}),1e3*t)})).catch((e=>{console.log(e);const t=this.$route.query.delay||5;console.log("delay",t),setTimeout((()=>{s.Z.replace({path:"/",query:{delay:t}})}),1e3*t)}))},onError(e){alert(e)},onSubmit(){this.arkoseToken||window.myEnforcement.run()}}},c=l,d=o(1),u=(0,d.Z)(c,n,r,!1,null,null,null),a=u.exports}}]);
//# sourceMappingURL=561-legacy.7504817a.js.map