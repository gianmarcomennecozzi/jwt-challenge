import Vue from 'vue'
import JWT from './components/jwt-form'

Vue.config.productionTip = false

if (document.getElementById("jwt-form")) {
  new Vue({
    render: h => h(JWT)
  }).$mount('#jwt-form')
}
