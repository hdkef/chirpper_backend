(window.webpackJsonp=window.webpackJsonp||[]).push([[11],{x5bZ:function(t,e,r){"use strict";r.r(e),r.d(e,"RegisterModule",function(){return C});var i=r("ofXK"),n=r("3Pt+"),o=r("LRne"),a=r("vkgz"),s=r("JIr8"),b=r("AytR"),d=r("AITh"),l=r("fXoL"),c=r("l7P3"),m=r("tk/3"),u=r("Wp6s"),p=r("kmnG"),g=r("qFsG"),f=r("bTqV"),h=r("tyNb");let w=(()=>{class t{constructor(t,e){this.store=t,this.http=e,this.verifSent=!1}ngOnInit(){this.registerForm=new n.d({Username:new n.b(null,[n.k.required]),Password:new n.b(null,[n.k.required]),Email:new n.b(null,[n.k.required,n.k.email]),Code:new n.b(null,[n.k.required])})}sendCode(){return this.verifSent=!0,this.http.post(`${b.a.api}${b.a.sendemailverroute}`,JSON.stringify({Email:this.registerForm.value.Email})).pipe(Object(a.a)(t=>{alert(t.MESSAGE)}),Object(s.a)(t=>(this.verifSent=!1,Object(o.a)(this.store.dispatch(new d.l({Info:"cannot send code"})))))).subscribe()}register(){let t=this.registerForm.value.Username,e=this.registerForm.value.Password,r=this.registerForm.value.Email;return this.http.post(`${b.a.api}${b.a.verifyemailverroute}`,JSON.stringify({Email:r,Code:this.registerForm.value.Code})).pipe(Object(a.a)(i=>{if(1!=i)return this.store.dispatch(new d.l({Info:"something wrong with email verification"}));this.store.dispatch(new d.j({Username:t,Password:e,Email:r}))})).subscribe()}}return t.\u0275fac=function(e){return new(e||t)(l.Mb(c.b),l.Mb(m.b))},t.\u0275cmp=l.Gb({type:t,selectors:[["app-register"]],decls:23,vars:3,consts:[["id","container"],["id","card"],["id","register-title"],[3,"formGroup","ngSubmit"],["id","register-input"],["matInput","","placeholder","Username","formControlName","Username"],["matInput","","placeholder","Password","type","password","formControlName","Password"],["matInput","","placeholder","Email","type","email","formControlName","Email"],["matInput","","placeholder","Code verification","type","text","formControlName","Code"],["mat-button","","type","button",2,"width","100%",3,"disabled","click"],["id","button-group"],["mat-raised-button","","color","accent","type","submit",3,"disabled"],["routerLink","/login",2,"text-decoration","none"],["mat-stroked-button","","color","primary"]],template:function(t,e){1&t&&(l.Rb(0,"div",0),l.Rb(1,"mat-card",1),l.Rb(2,"mat-card-title",2),l.kc(3,"Register"),l.Qb(),l.Rb(4,"form",3),l.Yb("ngSubmit",function(){return e.register()}),l.Rb(5,"div",4),l.Rb(6,"mat-form-field"),l.Nb(7,"input",5),l.Qb(),l.Rb(8,"mat-form-field"),l.Nb(9,"input",6),l.Qb(),l.Rb(10,"mat-form-field"),l.Nb(11,"input",7),l.Qb(),l.Rb(12,"mat-form-field"),l.Nb(13,"input",8),l.Qb(),l.Rb(14,"div"),l.Rb(15,"button",9),l.Yb("click",function(){return e.sendCode()}),l.kc(16,"send Code"),l.Qb(),l.Qb(),l.Qb(),l.Rb(17,"div",10),l.Rb(18,"button",11),l.kc(19,"register"),l.Qb(),l.Rb(20,"a",12),l.Rb(21,"button",13),l.kc(22,"login"),l.Qb(),l.Qb(),l.Qb(),l.Qb(),l.Qb(),l.Qb()),2&t&&(l.Cb(4),l.dc("formGroup",e.registerForm),l.Cb(11),l.dc("disabled",e.registerForm.get("Email").invalid||e.verifSent),l.Cb(3),l.dc("disabled",e.registerForm.invalid))},directives:[u.a,u.h,n.l,n.h,n.e,p.b,g.a,n.a,n.g,n.c,f.a,h.c],styles:["#container[_ngcontent-%COMP%]{position:absolute;height:100%;width:100%;background-color:indigo}#register-title[_ngcontent-%COMP%]{margin-bottom:15px}#register-input[_ngcontent-%COMP%]   mat-form-field[_ngcontent-%COMP%]{display:block}#card[_ngcontent-%COMP%]{position:relative;margin:20vh auto 0;width:200px;padding:50px}#button-group[_ngcontent-%COMP%]   button[_ngcontent-%COMP%]{display:block}#button-group[_ngcontent-%COMP%]{display:flex;margin-top:15px;justify-content:space-between}"]}),t})();var v=r("PCNd");let C=(()=>{class t{}return t.\u0275fac=function(e){return new(e||t)},t.\u0275mod=l.Kb({type:t}),t.\u0275inj=l.Jb({imports:[[i.b,h.d.forChild([{path:"",component:w}]),v.a]]}),t})()}}]);