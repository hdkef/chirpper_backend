(window.webpackJsonp=window.webpackJsonp||[]).push([[6],{ErIi:function(t,e,n){"use strict";n.r(e),n.d(e,"CommentModule",function(){return v});var r=n("ofXK"),o=n("7qN/"),c=n("3Pt+"),i=n("fXoL"),m=n("l7P3"),a=n("tyNb"),s=n("2dwN"),d=n("Wp6s");function b(t,e){if(1&t&&(i.Rb(0,"div",3),i.Nb(1,"img",4),i.Qb()),2&t){const t=i.ac();i.Cb(1),i.dc("src",t.chirp.ImageURL,i.hc)}}const p=function(t){return{"background-image":t}};let u=(()=>{class t{constructor(t){this.router=t}ngOnInit(){}goToProfile(){this.router.navigate(["profile"],{queryParams:{ID:this.chirp.ID}})}}return t.\u0275fac=function(e){return new(e||t)(i.Mb(a.b))},t.\u0275cmp=i.Gb({type:t,selectors:[["app-chirp"]],inputs:{chirp:"chirp"},decls:11,vars:7,consts:[["mat-card-avatar","",2,"background-size","cover",3,"ngStyle"],[2,"cursor","pointer",3,"click"],["id","imgheader",4,"ngIf"],["id","imgheader"],[3,"src"]],template:function(t,e){1&t&&(i.Rb(0,"mat-card"),i.Rb(1,"mat-card-header"),i.Nb(2,"div",0),i.Rb(3,"mat-card-title",1),i.Yb("click",function(){return e.goToProfile()}),i.kc(4),i.Qb(),i.Rb(5,"mat-card-subtitle"),i.kc(6),i.Qb(),i.Qb(),i.jc(7,b,2,1,"div",2),i.Rb(8,"mat-card-content"),i.Rb(9,"p"),i.kc(10),i.Qb(),i.Qb(),i.Qb()),2&t&&(i.Cb(2),i.dc("ngStyle",i.ec(5,p,"url("+e.chirp.AvatarURL+")")),i.Cb(2),i.lc(e.chirp.Username),i.Cb(2),i.lc(e.chirp.Date),i.Cb(1),i.dc("ngIf",e.chirp.ImageURL),i.Cb(3),i.mc(" ",e.chirp.Text," "))},directives:[d.a,d.e,d.c,r.j,d.h,d.g,r.i,d.d],styles:["mat-card[_ngcontent-%COMP%]{background-color:#f5f5f5;margin-top:15px}#imgheader[_ngcontent-%COMP%]{display:flex;justify-content:center;margin-bottom:25px}#imgheader[_ngcontent-%COMP%]   img[_ngcontent-%COMP%]{width:100%;overflow-y:hidden}mat-card-content[_ngcontent-%COMP%]{padding-top:10px}"]}),t})();var h=n("kmnG"),l=n("qFsG"),g=n("bTqV");function f(t,e){1&t&&i.Nb(0,"app-chirp",1),2&t&&i.dc("chirp",e.$implicit)}let C=(()=>{class t{constructor(t,e,n){this.store=t,this.route=e,this.wscomment=n}ngOnDestroy(){this.commentSubs&&this.commentSubs.unsubscribe(),this.store.dispatch(new o.l({}))}ngOnInit(){this.commentSubs=this.store.select("endpoints").subscribe(t=>{this.commentHeader=t.commentHeader,this.comment=t.comment}),this.PostID=this.route.snapshot.queryParamMap.get("PostID"),this.store.dispatch(new o.o(this.PostID)),this.commentForm=new c.d({Text:new c.b(null,c.k.required)})}sendComment(){this.wscomment.sendMsgPayload(this.commentForm.value.Text),this.commentForm.setValue({Text:null}),this.commentForm.markAsPristine(),this.commentForm.markAsUntouched(),this.commentForm.controls.Text.setErrors(null)}}return t.\u0275fac=function(e){return new(e||t)(i.Mb(m.b),i.Mb(a.a),i.Mb(s.a))},t.\u0275cmp=i.Gb({type:t,selectors:[["app-comment"]],decls:11,vars:4,consts:[["id","commentHeader"],[3,"chirp"],[3,"chirp",4,"ngFor","ngForOf"],[3,"formGroup","ngSubmit"],["id","txt"],["appearance","fill"],["matInput","","type","text","placeholder","add a comment here...","formControlName","Text"],["id","btn-group"],["mat-raised-button","","color","primary",3,"disabled"]],template:function(t,e){1&t&&(i.Rb(0,"div",0),i.Nb(1,"app-chirp",1),i.Qb(),i.jc(2,f,1,1,"app-chirp",2),i.Rb(3,"div"),i.Rb(4,"form",3),i.Yb("ngSubmit",function(){return e.sendComment()}),i.Rb(5,"div",4),i.Rb(6,"mat-form-field",5),i.Nb(7,"textarea",6),i.Qb(),i.Qb(),i.Rb(8,"div",7),i.Rb(9,"button",8),i.kc(10,"send"),i.Qb(),i.Qb(),i.Qb(),i.Qb()),2&t&&(i.Cb(1),i.dc("chirp",e.commentHeader),i.Cb(1),i.dc("ngForOf",e.comment),i.Cb(2),i.dc("formGroup",e.commentForm),i.Cb(5),i.dc("disabled",e.commentForm.invalid))},directives:[u,r.h,c.l,c.h,c.e,h.b,l.a,c.a,c.g,c.c,g.a],styles:["#commentHeader[_ngcontent-%COMP%]{margin-bottom:20px}#btn-group[_ngcontent-%COMP%]   button[_ngcontent-%COMP%], #txt[_ngcontent-%COMP%]   mat-form-field[_ngcontent-%COMP%], #txt[_ngcontent-%COMP%]   mat-form-field[_ngcontent-%COMP%]   textarea[_ngcontent-%COMP%]{width:100%}"]}),t})();var P=n("PCNd");let v=(()=>{class t{}return t.\u0275fac=function(e){return new(e||t)},t.\u0275mod=i.Kb({type:t}),t.\u0275inj=i.Jb({imports:[[r.b,a.d.forChild([{path:"",component:C}]),P.a]]}),t})()}}]);