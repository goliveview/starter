import {Application} from "@hotwired/stimulus"
import {definitionsFromContext} from "@hotwired/stimulus-webpack-helpers"
import {GoliveviewController} from 'goliveview-stimulus';
import "./styles.scss";

const application = Application.start()
const context = require.context("./controllers", true, /\.js$/)
application.load(definitionsFromContext(context))
application.register("glv", GoliveviewController)
window.Stimulus = application
