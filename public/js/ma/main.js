import $ from 'jquery'
import Bootstrap from 'bootstrap/dist/js/npm'
import _Initializer from './init/_Initializer'
import BootstrapExtend from './BootstrapExtend'

$(document).ready(() => {
    // switch initializer by URL path
    _Initializer.getInitializer().init()
})