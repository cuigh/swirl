var Swirl;
(function (Swirl) {
    var Core;
    (function (Core) {
        class Modal {
            constructor(modal) {
                this.$modal = $(modal);
                this.find(".modal-background, .modal-close, .modal-card-head .delete, .modal-card-foot .dismiss").click(e => this.close());
            }
            static initialize() {
                if (!Modal.initialized) {
                    $('.modal-trigger').click(function () {
                        let target = $(this).data('target');
                        let dlg = new Modal("#" + target);
                        dlg.show();
                    });
                    $(document).on('keyup', function (e) {
                        if (e.keyCode === 27) {
                            Modal.active && Modal.active.close();
                        }
                    });
                    Modal.initialized = true;
                }
            }
            static close() {
                Modal.active && Modal.active.close();
            }
            static current() {
                return Modal.active;
            }
            static alert(content, title, callback) {
                title = title || "Prompt";
                let $elem = $(`<div class="modal">
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">${title}</p>
                <button class="delete"></button>
              </header>
              <section class="modal-card-body">${content}</section>
              <footer class="modal-card-foot">
                <button class="button is-primary">OK</button>
              </footer>
            </div>
          </div>`).appendTo("body");
                let dlg = new Modal($elem);
                callback = callback || function (dlg, e) { dlg.close(); };
                $elem.find(".modal-card-foot>button:first").click(e => callback(dlg, e));
                dlg.deletable = true;
                dlg.show();
                return dlg;
            }
            static confirm(content, title, callback) {
                title = title || "Confirm";
                let $elem = $(`<div class="modal">
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">${title}</p>
                <button class="delete"></button>
              </header>
              <section class="modal-card-body">${content}</section>
              <footer class="modal-card-foot">
                <button class="button is-primary">Confirm</button>
                <button class="button dismiss">Cancel</button>
              </footer>
            </div>
          </div>`).appendTo("body");
                let dlg = new Modal($elem);
                if (callback) {
                    $elem.find(".modal-card-foot>button:first").click(e => callback(dlg, e));
                }
                dlg.deletable = true;
                dlg.show();
                return dlg;
            }
            show() {
                Modal.active && Modal.active.close();
                $('html').addClass('is-clipped');
                this.$modal.addClass('is-active').focus();
                this.error();
                Modal.active = this;
            }
            close() {
                $('html').removeClass('is-clipped');
                if (this.deletable) {
                    this.$modal.remove();
                }
                else {
                    this.$modal.removeClass('is-active');
                }
                Modal.active = null;
            }
            error(msg) {
                let $error = this.find(".modal-card-error");
                if (msg) {
                    if (!$error.length) {
                        $error = $('<section class="modal-card-error"></section>').insertAfter(this.find(".modal-card-body"));
                    }
                    $error.html(msg).show();
                }
                else {
                    $error.hide();
                }
            }
            find(selector) {
                if (typeof selector === "string") {
                    return this.$modal.find(selector);
                }
                else if (selector instanceof Element) {
                    return this.$modal.find(selector);
                }
                return this.$modal.find(selector);
            }
        }
        Core.Modal = Modal;
        class Notification {
            constructor(options) {
                this.options = $.extend({}, {
                    type: "primary",
                    time: 5,
                }, options);
            }
            static show(type, msg, time) {
                let n = new Notification({
                    type: type,
                    message: msg,
                    time: time,
                });
                n.show();
                return n;
            }
            show() {
                this.$elem = $(`<div class="notification is-${this.options.type} has-text-centered is-marginless is-radiusless is-top" style="display:none">${this.options.message}</div>`).appendTo("body").fadeIn(250);
                if (this.options.time > 0) {
                    setTimeout(() => this.hide(), this.options.time * 1000);
                }
            }
            hide() {
                this.$elem.fadeTo("slow", 0.01, () => this.$elem.remove());
            }
        }
        Core.Notification = Notification;
        class Tab {
            constructor(tab, content) {
                this.$tab = $(tab);
                this.$content = $(content);
                this.active = this.$tab.find("li.is-active").index();
                if (this.active == -1) {
                    this.select(0);
                }
                this.$tab.find("li").click(e => {
                    let $li = $(e.target).closest("li");
                    if (!$li.hasClass("is-active")) {
                        this.select($li.index());
                    }
                });
            }
            static initialize() {
                if (!Tab.initialized) {
                    $('.tabs').each((i, elem) => {
                        let target = $(elem).data("target");
                        new Tab(elem, "#" + target);
                    });
                    Tab.initialized = true;
                }
            }
            select(index) {
                if (this.active != index) {
                    this.$tab.find(`li.is-active`).removeClass("is-active");
                    this.$tab.find(`li:eq(${index})`).addClass("is-active");
                    this.$content.children(":visible").hide();
                    $(this.$content.children()[index]).show();
                    this.active = index;
                }
            }
        }
        Core.Tab = Tab;
        class FilterBox {
            constructor(elem, callback, timeout) {
                this.$elem = $(elem);
                this.$elem.keyup(() => {
                    if (this.timer > 0) {
                        clearTimeout(this.timer);
                    }
                    let text = this.$elem.val().toLowerCase();
                    this.timer = setTimeout(() => callback(text), timeout || 500);
                });
            }
        }
        Core.FilterBox = FilterBox;
    })(Core = Swirl.Core || (Swirl.Core = {}));
})(Swirl || (Swirl = {}));
/*!
 * Swirl Ajax Library v1.0.0
 * Copyright 2017 cuigh. All rights reserved.
 *
 * @author cuigh(noname@live.com)
 */
var Swirl;
(function (Swirl) {
    var Core;
    (function (Core) {
        class AjaxResult {
        }
        Core.AjaxResult = AjaxResult;
        class AjaxOptions {
            constructor() {
                this.timeout = 30000;
                this.async = true;
                this.encoder = "json";
            }
            static getDefaultOptions() {
                return AjaxOptions.defaultOptions;
            }
            static setDefaultOptions(options) {
                if (options) {
                    AjaxOptions.defaultOptions = options;
                }
            }
        }
        AjaxOptions.defaultOptions = new AjaxOptions();
        Core.AjaxOptions = AjaxOptions;
        let AjaxMethod;
        (function (AjaxMethod) {
            AjaxMethod[AjaxMethod["GET"] = 0] = "GET";
            AjaxMethod[AjaxMethod["POST"] = 1] = "POST";
            AjaxMethod[AjaxMethod["PUT"] = 2] = "PUT";
            AjaxMethod[AjaxMethod["DELETE"] = 3] = "DELETE";
            AjaxMethod[AjaxMethod["HEAD"] = 4] = "HEAD";
            AjaxMethod[AjaxMethod["TRACE"] = 5] = "TRACE";
            AjaxMethod[AjaxMethod["OPTIONS"] = 6] = "OPTIONS";
            AjaxMethod[AjaxMethod["CONNECT"] = 7] = "CONNECT";
            AjaxMethod[AjaxMethod["PATCH"] = 8] = "PATCH";
        })(AjaxMethod = Core.AjaxMethod || (Core.AjaxMethod = {}));
        class AjaxRequest {
            constructor(url, method, data) {
                this.options = $.extend({
                    url: url,
                    method: method,
                    data: data,
                    preHandler: AjaxRequest.preHandler,
                    postHandler: AjaxRequest.postHandler,
                    errorHandler: AjaxRequest.errorHandler,
                }, AjaxOptions.getDefaultOptions());
            }
            preHandler(handler) {
                this.options.preHandler = handler;
                return this;
            }
            postHandler(handler) {
                this.options.postHandler = handler;
                return this;
            }
            errorHandler(handler) {
                this.options.errorHandler = handler;
                return this;
            }
            timeout(timeout) {
                this.options.timeout = timeout;
                return this;
            }
            async(async) {
                this.options.async = async;
                return this;
            }
            trigger(elem) {
                this.options.trigger = elem;
                return this;
            }
            json(callback) {
                return this.result("json", callback);
            }
            text(callback) {
                return this.result("text", callback);
            }
            html(callback) {
                return this.result("html", callback);
            }
            result(dataType, callback) {
                this.options.dataType = dataType;
                this.options.preHandler && this.options.preHandler(this.options);
                let settings = this.buildSettings();
                if (callback) {
                    $.ajax(settings).done(callback).always(() => {
                        this.options.postHandler && this.options.postHandler(this.options);
                    }).fail((xhr, textStatus, error) => {
                        this.options.errorHandler && this.options.errorHandler(xhr, textStatus, error);
                    });
                }
                else {
                    return new Promise((resolve, _) => {
                        $.ajax(settings).done((r) => {
                            resolve(r);
                        }).always(() => {
                            AjaxRequest.postHandler && AjaxRequest.postHandler(this.options);
                        }).fail((xhr, textStatus, error) => {
                            AjaxRequest.errorHandler && AjaxRequest.errorHandler(xhr, textStatus, error);
                        });
                    });
                }
            }
            buildSettings() {
                return {
                    url: this.options.url,
                    method: AjaxMethod[this.options.method],
                    data: this.options.data,
                    dataType: this.options.dataType,
                    timeout: this.options.timeout,
                    async: this.options.async,
                };
            }
        }
        AjaxRequest.preHandler = options => {
            options.trigger && $(options.trigger).prop("disabled", true);
        };
        AjaxRequest.postHandler = options => {
            options.trigger && $(options.trigger).prop("disabled", false);
        };
        AjaxRequest.errorHandler = (xhr, status, err) => {
            let msg;
            if (xhr.responseJSON) {
                let err = xhr.responseJSON;
                msg = err.message;
                if (err.code) {
                    msg += `(${err.code})`;
                }
            }
            else if (xhr.status >= 400) {
                msg = xhr.responseText || err || status;
            }
            else {
                return;
            }
            Core.Notification.show("danger", `AJAX: ${msg}`);
        };
        Core.AjaxRequest = AjaxRequest;
        class AjaxGetRequest extends AjaxRequest {
            jsonp(callback) {
                return this.result("jsonp", callback);
            }
        }
        Core.AjaxGetRequest = AjaxGetRequest;
        class AjaxPostRequest extends AjaxRequest {
            encoder(encoder) {
                this.options.encoder = encoder;
                return this;
            }
            buildSettings() {
                let settings = super.buildSettings();
                switch (this.options.encoder) {
                    case "none":
                        settings.contentType = false;
                        settings.processData = false;
                        break;
                    case "json":
                        settings.contentType = "application/json; charset=UTF-8";
                        settings.data = JSON.stringify(this.options.data);
                        break;
                    case "form":
                        settings.contentType = "application/x-www-form-urlencoded; charset=UTF-8";
                        break;
                }
                return settings;
            }
        }
        Core.AjaxPostRequest = AjaxPostRequest;
        class Ajax {
            constructor() {
            }
            static get(url, args) {
                return new AjaxGetRequest(url, AjaxMethod.GET, args);
            }
            static post(url, data) {
                return new AjaxPostRequest(url, AjaxMethod.POST, data);
            }
        }
        Core.Ajax = Ajax;
    })(Core = Swirl.Core || (Swirl.Core = {}));
})(Swirl || (Swirl = {}));
let $ajax = Swirl.Core.Ajax;
var Swirl;
(function (Swirl) {
    var Core;
    (function (Core) {
        class NativeRule {
            validate($form, $input, arg) {
                let el = $input[0];
                return { ok: el.checkValidity ? el.checkValidity() : true };
            }
        }
        class RequiredRule {
            validate($form, $input, arg) {
                return { ok: $.trim($input.val()).length > 0 };
            }
        }
        class CheckedRule {
            validate($form, $input, arg) {
                let count = parseInt(arg);
                let siblings = $form.find(`:input:checked[name='${$input.attr("name")}']`);
                return { ok: siblings.length >= count };
            }
        }
        class EmailRule {
            validate($form, $input, arg) {
                const regex = /^((([a-z]|\d|[!#\$%&'\*\+\-\/=\?\^_`{\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+(\.([a-z]|\d|[!#\$%&'\*\+\-\/=\?\^_`{\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+)*)|((\x22)((((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(([\x01-\x08\x0b\x0c\x0e-\x1f\x7f]|\x21|[\x23-\x5b]|[\x5d-\x7e]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(\\([\x01-\x09\x0b\x0c\x0d-\x7f]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))))*(((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(\x22)))@((([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.)+(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))$/i;
                let value = $.trim($input.val());
                return { ok: !value || regex.test(value) };
            }
        }
        class UrlRule {
            validate($form, $input, arg) {
                const regex = /^(https?|ftp):\/\/(((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:)*@)?(((\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5]))|((([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.)+(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.?)(:\d*)?)(\/((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)+(\/(([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)*)*)?)?(\?((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)|[\uE000-\uF8FF]|\/|\?)*)?(\#((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)|\/|\?)*)?$/i;
                let value = $.trim($input.val());
                return { ok: !value || regex.test(value) };
            }
        }
        class IPRule {
            validate($form, $input, arg) {
                const regex = /^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[0-9]{1,2})$/i;
                let value = $.trim($input.val());
                return { ok: !value || regex.test(value) };
            }
        }
        class MatchValidator {
            validate($form, $input, arg) {
                return { ok: $input.val() == $('#' + arg).val() };
            }
        }
        class LengthRule {
            validate($form, $input, arg) {
                let r = { ok: true };
                if (arg) {
                    let len = this.getLength($.trim($input.val()));
                    let args = arg.split('~');
                    if (args.length == 1) {
                        if ($.isNumeric(args[0])) {
                            r.ok = len >= parseInt(args[0]);
                        }
                    }
                    else {
                        if ($.isNumeric(args[0]) && $.isNumeric(args[1])) {
                            r.ok = len >= parseInt(args[0]) && len <= parseInt(args[1]);
                        }
                    }
                }
                return r;
            }
            getLength(value) {
                return value.length;
            }
        }
        class WidthRule extends LengthRule {
            getLength(value) {
                let doubleByteChars = value.match(/[^\x00-\xff]/ig);
                return value.length + (doubleByteChars == null ? 0 : doubleByteChars.length);
            }
        }
        class IntegerRule {
            validate($form, $input, arg) {
                const regex = /^\d*$/;
                return { ok: regex.test($.trim($input.val())) };
            }
        }
        class RegexRule {
            validate($form, $input, arg) {
                let regex = new RegExp(arg);
                let value = $.trim($input.val());
                return { ok: !value || regex.test(value) };
            }
        }
        class RemoteRule {
            validate($form, $input, arg) {
                if (!arg) {
                    throw new Error("服务器验证地址未设置");
                }
                let value = $.trim($input.val());
                let r = { ok: false };
                $ajax.post(arg, { value: value }).encoder("form").async(false).json(result => {
                    r.ok = !result.error;
                    r.error = result.error;
                });
                return r;
            }
        }
        class Validator {
            constructor(elem, options) {
                this.form = $(elem);
                this.options = options;
                if (this.form.is("form")) {
                    this.form.attr("novalidate", "true");
                }
                this.form.on("click", ':radio[data-v-rule],:checkbox[data-v-rule]', this.checkValue.bind(this));
                this.form.on("change", 'select[data-v-rule],input[type="file"][data-v-rule]', this.checkValue.bind(this));
                this.form.on("blur", ':input[data-v-rule]:not(select,:radio,:checkbox,:file)', this.checkValue.bind(this));
            }
            checkValue(e) {
                let $input = $(e.target);
                let result = this.validateInput($input);
                Validator.mark($input, result);
            }
            static bind(elem, options) {
                let v = $(elem).data("validator");
                if (!v) {
                    v = new Validator(elem, options);
                    $(elem).data("validator", v);
                }
                return v;
            }
            validate() {
                let results = [];
                this.form.find(Validator.selector).each((i, el) => {
                    let $input = $(el);
                    let result = this.validateInput($input);
                    if (result != null) {
                        results.push(result);
                    }
                    Validator.mark($input, result);
                });
                return results;
            }
            reset() {
                this.form.find(Validator.selector).each((i, el) => {
                    let $input = $(el);
                    Validator.marker.reset($input);
                });
            }
            static register(name, rule, msg) {
                this.rules[name] = rule;
                this.messages[name] = msg;
            }
            static setMessage(name, msg) {
                this.messages[name] = msg;
            }
            static setMarker(marker) {
                this.marker = marker;
            }
            validateInput($input) {
                let errors = [];
                let rules = ($input.data('v-rule') || 'native').split(';');
                rules.forEach(name => {
                    let rule = Validator.rules[name];
                    if (rule) {
                        let arg = $input.data("v-arg-" + name);
                        let r = rule.validate(this.form, $input, arg);
                        if (!r.ok) {
                            errors.push(r.error || Validator.getMessge($input, name));
                        }
                    }
                });
                return (errors.length == 0) ? null : {
                    input: $input,
                    errors: errors,
                };
            }
            static mark($input, result) {
                if (Validator.marker != null) {
                    if (result) {
                        Validator.marker.setError($input, result.errors);
                    }
                    else {
                        Validator.marker.clearError($input);
                    }
                }
            }
            static getMessge($input, rule) {
                if (rule == 'native')
                    return $input[0].validationMessage;
                else {
                    let msg = $input.data('v-msg-' + rule);
                    if (!msg)
                        msg = this.messages[rule];
                    return msg;
                }
            }
        }
        Validator.selector = ':input[data-v-rule]:not(:submit,:button,:reset,:image,:disabled)';
        Validator.messages = {
            "required": "This field is required",
            "checked": "Number of checked items is invalid",
            "email": "Please input a valid email address",
            "match": "Input confirmation doesn't match",
            "length": "The length of the field does not meet the requirements",
            "width": "The width of the field does not meet the requirements",
            "url": "Please input a valid url",
            "ip": "Please input a valid IPV4 address",
            "integer": "Please input an integer",
            "regex": "Input is invalid",
            "remote": "Input is invalid",
        };
        Validator.rules = {
            "native": new NativeRule(),
            "required": new RequiredRule(),
            "checked": new CheckedRule(),
            "email": new EmailRule(),
            "match": new MatchValidator(),
            "length": new LengthRule(),
            "width": new WidthRule(),
            "url": new UrlRule(),
            "ip": new IPRule(),
            "integer": new IntegerRule(),
            "regex": new RegexRule(),
            "remote": new RemoteRule(),
        };
        Core.Validator = Validator;
    })(Core = Swirl.Core || (Swirl.Core = {}));
})(Swirl || (Swirl = {}));
/*!
 * Swirl Form Library v1.0.0
 * Copyright 2017 cuigh. All rights reserved.
 * see also: https://github.com/A1rPun/transForm.js
 *
 * @author cuigh(noname@live.com)
 */
var Swirl;
(function (Swirl) {
    var Core;
    (function (Core) {
        class FormOptions {
            constructor() {
                this.delimiter = ".";
                this.skipDisabled = true;
                this.skipReadOnly = false;
                this.skipEmpty = false;
                this.useIdOnEmptyName = true;
                this.triggerChange = false;
            }
        }
        Core.FormOptions = FormOptions;
        class Form {
            constructor(elem, options) {
                this.form = $(elem);
                this.options = $.extend(new FormOptions(), options);
                this.validator = Core.Validator.bind(this.form);
            }
            reset() {
                this.form.get(0).reset();
                if (this.validator) {
                    this.validator.reset();
                }
            }
            clear() {
                let inputs = this.getFields();
                inputs.each((i, input) => {
                    this.clearInput(input);
                });
                if (this.validator) {
                    this.validator.reset();
                }
            }
            submit(url) {
                let data = this.serialize();
                return Core.Ajax.post(url || this.form.attr("action"), data);
            }
            validate() {
                if (!this.validator) {
                    return true;
                }
                return this.validator.validate().length == 0;
            }
            serialize(nodeCallback) {
                let result = {}, inputs = this.getFields();
                for (let i = 0, l = inputs.length; i < l; i++) {
                    let input = inputs[i], key = input.name || this.options.useIdOnEmptyName && input.id;
                    if (!key)
                        continue;
                    let entry = null;
                    if (nodeCallback)
                        entry = nodeCallback(input);
                    if (!entry)
                        entry = this.getEntryFromInput(input, key);
                    if (typeof entry.value === 'undefined' || entry.value === null
                        || (this.options.skipEmpty && (!entry.value || (this.isArray(entry.value) && !entry.value.length))))
                        continue;
                    this.saveEntryToResult(result, entry, input, this.options.delimiter);
                }
                return result;
            }
            deserialize(data, nodeCallback) {
                let inputs = this.getFields();
                for (let i = 0, l = inputs.length; i < l; i++) {
                    let input = inputs[i], key = input.name || this.options.useIdOnEmptyName && input.id, value = this.getFieldValue(key, data);
                    if (typeof value === 'undefined' || value === null) {
                        this.clearInput(input);
                        continue;
                    }
                    let mutated = nodeCallback && nodeCallback(input, value);
                    if (!mutated)
                        this.setValueToInput(input, value);
                }
            }
            static automate() {
                $('form[data-form]').each(function (i, elem) {
                    let $form = $(elem);
                    let form = new Form($form);
                    let type = $form.data("form");
                    if (type == "form") {
                        $form.submit(e => { return form.validate(); });
                        return;
                    }
                    $form.submit(function () {
                        if (!form.validate()) {
                            return false;
                        }
                        let request = form.submit($form.attr("action")).trigger($form.find('button[type="submit"]'));
                        if (type == "ajax-form") {
                            request.encoder("form");
                        }
                        request.json((r) => {
                            if (r.success) {
                                let url = r.url || $form.data("url");
                                if (url) {
                                    if (url === "-") {
                                        location.reload();
                                    }
                                    else {
                                        location.href = url;
                                    }
                                }
                                else {
                                    let msg = r.message || $form.data("message");
                                    Core.Notification.show("info", `SUCCESS: ${msg}`, 3);
                                }
                            }
                            else {
                                let msg = r.message;
                                if (r.code) {
                                    msg += `({r.code})`;
                                }
                                Core.Notification.show("danger", `FAILED: ${msg}`);
                            }
                        });
                        return false;
                    });
                });
            }
            getEntryFromInput(input, key) {
                let nodeType = input.type && input.type.toLowerCase(), entry = { name: key }, dataType = $(input).data("type");
                switch (nodeType) {
                    case 'radio':
                        if (input.checked)
                            entry.value = input.value === 'on' ? true : input.value;
                        break;
                    case 'checkbox':
                        entry.value = input.checked ? (input.value === 'on' ? true : input.value) : false;
                        break;
                    case 'select-multiple':
                        entry.value = [];
                        for (let i = 0, l = input.options.length; i < l; i++)
                            if (input.options[i].selected)
                                entry.value.push(input.options[i].value);
                        break;
                    case 'file':
                        entry.value = input.value.split('\\').pop();
                        break;
                    case 'button':
                    case 'submit':
                    case 'reset':
                        break;
                    default:
                        entry.value = input.value;
                }
                if (entry.value != null) {
                    switch (dataType) {
                        case "integer":
                            entry.value = parseInt(entry.value);
                            break;
                        case "float":
                            entry.value = parseFloat(entry.value);
                            break;
                        case "bool":
                            entry.value = (entry.value === "true") || (entry.value === "1");
                            break;
                    }
                }
                return entry;
            }
            saveEntryToResult(parent, entry, input, delimiter) {
                if (/\[]$/.test(entry.name) && !entry.value)
                    return;
                let parts = this.parseString(entry.name, delimiter);
                for (let i = 0, l = parts.length; i < l; i++) {
                    let part = parts[i];
                    if (i === l - 1) {
                        parent[part] = entry.value;
                    }
                    else {
                        let index = parts[i + 1];
                        if (!index || $.isNumeric(index)) {
                            if (!this.isArray(parent[part]))
                                parent[part] = [];
                            if (i === l - 2) {
                                parent[part].push(entry.value);
                            }
                            else {
                                if (!this.isObject(parent[part][index]))
                                    parent[part][index] = {};
                                parent = parent[part][index];
                            }
                            i++;
                        }
                        else {
                            if (!this.isObject(parent[part]))
                                parent[part] = {};
                            parent = parent[part];
                        }
                    }
                }
            }
            clearInput(input) {
                let nodeType = input.type && input.type.toLowerCase();
                switch (nodeType) {
                    case 'select-one':
                        input.selectedIndex = 0;
                        break;
                    case 'select-multiple':
                        for (let i = input.options.length; i--;)
                            input.options[i].selected = false;
                        break;
                    case 'radio':
                    case 'checkbox':
                        if (input.checked)
                            input.checked = false;
                        break;
                    case 'button':
                    case 'submit':
                    case 'reset':
                    case 'file':
                        break;
                    default:
                        input.value = '';
                }
                if (this.options.triggerChange) {
                    $(input).change();
                }
            }
            parseString(str, delimiter) {
                let result = [], split = str.split(delimiter), len = split.length;
                for (let i = 0; i < len; i++) {
                    let s = split[i].split('['), l = s.length;
                    for (let j = 0; j < l; j++) {
                        let key = s[j];
                        if (!key) {
                            if (j === 0)
                                continue;
                            if (j !== l - 1)
                                throw new Error(`Undefined key is not the last part of the name > ${str}`);
                        }
                        if (key && key[key.length - 1] === ']')
                            key = key.slice(0, -1);
                        result.push(key);
                    }
                }
                return result;
            }
            getFields() {
                let inputs = this.form.find("input,select,textarea").filter(':not([data-form-ignore="true"])');
                if (this.options.skipDisabled)
                    inputs = inputs.filter(":not([disabled])");
                if (this.options.skipReadOnly)
                    inputs = inputs.filter(":not([readonly])");
                return inputs;
            }
            getFieldValue(key, ref) {
                if (!key || !ref)
                    return;
                let parts = this.parseString(key, this.options.delimiter);
                for (let i = 0, l = parts.length; i < l; i++) {
                    let part = ref[parts[i]];
                    if (typeof part === 'undefined' || part === null)
                        return;
                    if (i === l - 1) {
                        return part;
                    }
                    else {
                        let index = parts[i + 1];
                        if (index === '') {
                            return part;
                        }
                        else if ($.isNumeric(index)) {
                            if (i === l - 2)
                                return part[index];
                            else
                                ref = part[index];
                            i++;
                        }
                        else {
                            ref = part;
                        }
                    }
                }
            }
            setValueToInput(input, value) {
                let nodeType = input.type && input.type.toLowerCase();
                switch (nodeType) {
                    case 'radio':
                        if (value == input.value)
                            input.checked = true;
                        break;
                    case 'checkbox':
                        input.checked = this.isArray(value)
                            ? this.contains(value, input.value)
                            : value === true || value == input.value;
                        break;
                    case 'select-multiple':
                        if (this.isArray(value))
                            for (let i = input.options.length; i--;)
                                input.options[i].selected = this.contains(value, input.options[i].value);
                        else
                            input.value = value;
                        break;
                    case 'button':
                    case 'submit':
                    case 'reset':
                    case 'file':
                        break;
                    default:
                        input.value = value;
                }
                if (this.options.triggerChange) {
                    $(input).change();
                }
            }
            contains(array, value) {
                for (let item of array) {
                    if (item == value)
                        return true;
                }
                return false;
            }
            isObject(obj) {
                return typeof obj === 'object';
            }
            isArray(arr) {
                return Array.isArray(arr);
            }
        }
        Core.Form = Form;
    })(Core = Swirl.Core || (Swirl.Core = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Core;
    (function (Core) {
        class Dispatcher {
            constructor(attr) {
                this.events = {};
                this.attr = attr || "action";
            }
            static bind(elem, event = "click") {
                return new Dispatcher().bind(elem, event);
            }
            on(action, handler) {
                this.events[action] = handler;
                return this;
            }
            off(action) {
                delete this.events[action];
                return this;
            }
            bind(elem, event = "click") {
                $(elem).on(event, this.handle.bind(this));
                return this;
            }
            handle(e) {
                let action = $(e.target).closest("[data-" + this.attr + "]").data(this.attr);
                if (action) {
                    let handler = this.events[action];
                    if (handler) {
                        e.stopPropagation();
                        handler(e);
                    }
                }
            }
        }
        Core.Dispatcher = Dispatcher;
    })(Core = Swirl.Core || (Swirl.Core = {}));
})(Swirl || (Swirl = {}));
/*!
 * Swirl Table Library v1.0.0
 * Copyright 2017 cuigh. All rights reserved.
 *
 * @author cuigh(noname@live.com)
 */
var Swirl;
(function (Swirl) {
    var Core;
    (function (Core) {
        class Table {
            constructor(table) {
                this.$table = $(table);
                this.dispatcher = Core.Dispatcher.bind(this.$table);
            }
            on(action, handler) {
                this.dispatcher.on(action, handler);
                return this;
            }
        }
        Core.Table = Table;
        class ListTable extends Table {
            constructor(table) {
                super(table);
                this.on("check-all", e => {
                    let checked = e.target.checked;
                    this.$table.find("tbody>tr").each((i, elem) => {
                        $(elem).find("td:first>:checkbox").prop("checked", checked);
                    });
                });
                this.on("check", () => {
                    let rows = this.$table.find("tbody>tr").length;
                    let checkedRows = this.selectedRows().length;
                    this.$table.find("thead>tr>th:first>:checkbox").prop("checked", checkedRows > 0 && rows == checkedRows);
                });
            }
            selectedRows() {
                return this.$table.find("tbody>tr").filter((i, elem) => {
                    let cb = $(elem).find("td:first>:checkbox");
                    return cb.prop("checked");
                });
            }
            selectedKeys() {
                let keys = [];
                this.$table.find("tbody>tr").each((i, elem) => {
                    let cb = $(elem).find("td:first>:checkbox");
                    if (cb.prop("checked")) {
                        keys.push(cb.val());
                    }
                });
                return keys;
            }
        }
        Core.ListTable = ListTable;
        class EditTable extends Table {
            constructor(elem) {
                super(elem);
                this.name = this.$table.data("name");
                this.alias = this.name.replace(".", "-");
                this.index = this.$table.find("tbody>tr").length;
                super.on("add-" + this.alias, this.addRow.bind(this)).on("delete-" + this.alias, OptionTable.deleteRow);
            }
            addRow() {
                this.$table.find("tbody").append(this.render());
                this.index++;
            }
            static deleteRow(e) {
                $(e.target).closest("tr").remove();
            }
        }
        Core.EditTable = EditTable;
        class OptionTable extends EditTable {
            constructor(elem) {
                super(elem);
            }
            render() {
                return `<tr>
            <td><input name="${this.name}s[${this.index}].name" class="input is-small" type="text"></td>
            <td><input name="${this.name}s[${this.index}].value" class="input is-small" type="text"></td>
            <td>
              <a class="button is-small is-danger is-outlined" data-action="delete-${this.alias}">
                <span class="icon is-small"><i class="far fa-trash-alt"></i></span>
              </a>            
            </td>
          </tr>`;
            }
        }
        Core.OptionTable = OptionTable;
    })(Core = Swirl.Core || (Swirl.Core = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Core;
    (function (Core) {
        class BulmaMarker {
            setError($input, errors) {
                let $field = this.getField($input);
                $input.removeClass('is-success').addClass('is-danger');
                let $errors = $field.find('div.errors');
                if (!$errors.length) {
                    $errors = $('<div class="errors"/>').appendTo($field);
                }
                $errors.empty().append($.map(errors, (err) => `<p class="help is-danger">${err}</p>`)).show();
            }
            clearError($input) {
                let $field = this.getField($input);
                $input.removeClass('is-danger').addClass('is-success');
                let $errors = $field.find("div.errors");
                $errors.empty().hide();
            }
            reset($input) {
                let $field = this.getField($input);
                $input.removeClass('is-danger is-success');
                let $errors = $field.find("div.errors");
                $errors.empty().hide();
            }
            getField($input) {
                let $field = $input.closest(".field");
                if ($field.hasClass("has-addons") || $field.hasClass("is-grouped")) {
                    $field = $field.parent();
                }
                return $field;
            }
        }
        $(() => {
            $('.navbar-burger').click(function () {
                let $el = $(this);
                let $target = $('#' + $el.data('target'));
                $el.toggleClass('is-active');
                $target.toggleClass('is-active');
            });
            Core.Modal.initialize();
            Core.Tab.initialize();
            Core.AjaxRequest.preHandler = opts => opts.trigger && $(opts.trigger).addClass("is-loading");
            Core.AjaxRequest.postHandler = opts => opts.trigger && $(opts.trigger).removeClass("is-loading");
            Core.Validator.setMarker(new BulmaMarker());
            Core.Form.automate();
        });
    })(Core = Swirl.Core || (Swirl.Core = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Core;
    (function (Core) {
        class ChartOptions {
            constructor() {
                this.type = "line";
                this.width = 12;
                this.height = 200;
            }
        }
        Core.ChartOptions = ChartOptions;
        class Chart {
            constructor(opts) {
                this.opts = $.extend(new ChartOptions(), opts);
                this.createElem();
            }
            createElem() {
                this.$elem = $(`<div class="column is-${this.opts.width}" data-name="${this.opts.name}">
      <div class="card">
        <header class="card-header">
          <a class="card-header-icon drag">
            <span class="icon">
              <i class="fas fa-bars has-text-grey-light" aria-hidden="true"></i>
            </span>
          </a>        
          <p class="card-header-title is-paddingless">${this.opts.title}</p>
          <a data-action="remove-chart" class="card-header-icon" aria-label="remove chart">
            <span class="icon">
              <i class="far fa-trash-alt has-text-danger" aria-hidden="true"></i>
            </span>
          </a>         
        </header>
        <div class="card-content" style="height: ${this.opts.height}px"></div>
      </div>
    </div>`);
            }
            init() {
                let opt = {
                    legend: {
                        x: 'right',
                    },
                    tooltip: {
                        trigger: 'axis',
                        axisPointer: {
                            animation: false
                        }
                    },
                    xAxis: {
                        type: 'time',
                        splitNumber: 10,
                        splitLine: { show: false },
                    },
                    yAxis: {
                        type: 'value',
                        splitLine: {
                            lineStyle: {
                                type: "dashed",
                            },
                        },
                        axisLabel: {
                            formatter: this.formatValue.bind(this),
                        },
                    },
                    color: this.getColors(),
                };
                this.config(opt);
                this.chart = echarts.init(this.$elem.find("div.card-content").get(0));
                this.chart.setOption(opt, true);
            }
            getElem() {
                return this.$elem;
            }
            getOptions() {
                return this.opts;
            }
            resize() {
                this.chart.resize();
            }
            config(opt) {
            }
            formatValue(value) {
                switch (this.opts.unit) {
                    case "percent:100":
                        return value.toFixed(1) + "%";
                    case "percent:1":
                        return (value * 100).toFixed(1) + "%";
                    case "time:ns":
                        return value + 'ns';
                    case "time:µs":
                        return value.toFixed(2) + 'µs';
                    case "time:ms":
                        return value.toFixed(2) + 'ms';
                    case "time:s":
                        if (value < 1) {
                            return (value * 1000).toFixed(0) + 'ms';
                        }
                        else {
                            return value.toFixed(2) + 's';
                        }
                    case "time:m":
                        return value.toFixed(2) + 'm';
                    case "time:h":
                        return value.toFixed(2) + 'h';
                    case "time:d":
                        return value.toFixed(2) + 'd';
                    case "size:bits":
                        value = value / 8;
                    case "size:bytes":
                        if (value < 1024) {
                            return value.toString() + 'B';
                        }
                        else if (value < 1048576) {
                            return (value / 1024).toFixed(2) + 'K';
                        }
                        else if (value < 1073741824) {
                            return (value / 1048576).toFixed(2) + 'M';
                        }
                        else {
                            return (value / 1073741824).toFixed(2) + 'G';
                        }
                    case "size:kilobytes":
                        return value.toFixed(2) + 'K';
                    case "size:megabytes":
                        return value.toFixed(2) + 'M';
                    case "size:gigabytes":
                        return value.toFixed(2) + 'G';
                    default:
                        return value % 1 === 0 ? value.toString() : value.toFixed(2);
                }
            }
            getColors() {
                let colors = [
                    '#45aaf2',
                    '#6574cd',
                    '#a55eea',
                    '#f66d9b',
                    '#cd201f',
                    '#fd9644',
                    '#f1c40f',
                    '#7bd235',
                    '#5eba00',
                    '#2bcbba',
                ];
                this.shuffle(colors);
                return colors;
            }
            shuffle(a) {
                let len = a.length;
                for (let i = 0; i < len - 1; i++) {
                    let index = Math.floor(Math.random() * (len - i));
                    let temp = a[index];
                    a[index] = a[len - i - 1];
                    a[len - i - 1] = temp;
                }
            }
        }
        Core.Chart = Chart;
        class GaugeChart extends Chart {
            constructor(opts) {
                super(opts);
            }
            config(opt) {
                $.extend(true, opt, {
                    grid: {
                        left: 0,
                        top: 20,
                        right: 0,
                        bottom: 0,
                    },
                    xAxis: [
                        {
                            show: false,
                        },
                    ],
                    yAxis: [
                        {
                            show: false,
                        },
                    ],
                });
            }
            setData(d) {
                this.chart.setOption({
                    series: [
                        {
                            type: 'gauge',
                            radius: '100%',
                            center: ["50%", "58%"],
                            max: d.value,
                            axisLabel: { show: false },
                            pointer: { show: false },
                            detail: {
                                offsetCenter: [0, 0],
                            },
                            data: [{ value: d.value }]
                        }
                    ]
                });
            }
        }
        Core.GaugeChart = GaugeChart;
        class VectorChart extends Chart {
            constructor(opts) {
                super(opts);
            }
            config(opt) {
                $.extend(true, opt, {
                    grid: {
                        left: 20,
                        top: 20,
                        right: 20,
                        bottom: 20,
                    },
                    legend: {
                        type: 'scroll',
                        orient: 'vertical',
                    },
                    tooltip: {
                        trigger: 'item',
                        formatter: (params, index) => {
                            return params.name + ": " + this.formatValue(params.value);
                        },
                    },
                    xAxis: [
                        {
                            show: false,
                        },
                    ],
                    yAxis: [
                        {
                            show: false,
                        },
                    ],
                    series: [{
                            type: this.opts.type,
                            radius: '80%',
                            center: ['40%', '50%'],
                        }],
                });
            }
            setData(d) {
                this.chart.setOption({
                    legend: {
                        data: d.legend,
                    },
                    series: [{
                            data: d.data,
                        }],
                });
            }
        }
        Core.VectorChart = VectorChart;
        class MatrixChart extends Chart {
            constructor(opts) {
                super(opts);
            }
            config(opt) {
                $.extend(true, opt, {
                    grid: {
                        left: 60,
                        top: 30,
                        right: 20,
                        bottom: 30,
                    },
                    tooltip: {
                        formatter: (params) => {
                            let html = params[0].axisValueLabel + '<br/>';
                            for (let i = 0; i < params.length; i++) {
                                html += params[i].marker + params[i].seriesName + ': ' + this.formatValue(params[i].value[1]) + '<br/>';
                            }
                            return html;
                        },
                    },
                    yAxis: {
                        max: 'dataMax',
                    },
                });
            }
            setData(d) {
                if (!d.series) {
                    return;
                }
                d.series.forEach((s) => {
                    s.type = this.opts.type;
                    s.areaStyle = {
                        opacity: 0.3,
                    };
                    s.smooth = true;
                    s.showSymbol = false;
                    s.lineStyle = {
                        normal: {
                            width: 1,
                        }
                    };
                });
                this.chart.setOption({
                    legend: {
                        data: d.legend,
                    },
                    series: d.series,
                });
            }
        }
        Core.MatrixChart = MatrixChart;
        class ChartFactory {
            static create(opts) {
                switch (opts.type) {
                    case "gauge":
                        return new GaugeChart(opts);
                    case "line":
                    case "bar":
                        return new MatrixChart(opts);
                    case "pie":
                        return new VectorChart(opts);
                }
                return null;
            }
        }
        Core.ChartFactory = ChartFactory;
        class ChartDashboardOptions {
            constructor() {
                this.period = 30;
                this.refreshInterval = 15;
            }
        }
        Core.ChartDashboardOptions = ChartDashboardOptions;
        class ChartDashboard {
            constructor(elem, charts, opts) {
                this.charts = [];
                this.opts = $.extend(new ChartDashboardOptions(), opts);
                this.$panel = $(elem);
                this.dlg = new ChartDialog(this);
                charts.forEach(opts => this.createGraph(opts));
                Core.Dispatcher.bind(this.$panel).on("remove-chart", e => {
                    let name = $(e.target).closest("div.column").data("name");
                    Core.Modal.confirm(`Are you sure to delete chart: <strong>${name}</strong>?`, "Delete chart", dlg => {
                        this.removeGraph(name);
                        dlg.close();
                    });
                });
                $(window).resize(e => {
                    $.each(this.charts, (i, g) => {
                        g.resize();
                    });
                });
                this.refresh();
            }
            refresh() {
                if (!this.timer) {
                    this.loadData();
                    if (this.opts.refreshInterval > 0) {
                        this.timer = setTimeout(this.refreshData.bind(this), this.opts.refreshInterval * 1000);
                    }
                }
            }
            refreshData() {
                this.loadData();
                if (this.opts.refreshInterval > 0) {
                    this.timer = setTimeout(this.refreshData.bind(this), this.opts.refreshInterval * 1000);
                }
            }
            stop() {
                clearTimeout(this.timer);
                this.timer = 0;
            }
            setPeriod(period) {
                this.opts.period = period;
                this.loadData();
            }
            addGraph(opts) {
                this.createGraph(opts);
                this.loadData();
            }
            createGraph(opts) {
                for (let i = 0; i < this.charts.length; i++) {
                    let chart = this.charts[i];
                    if (chart.getOptions().name === opts.name) {
                        return;
                    }
                }
                let chart = ChartFactory.create(opts);
                if (chart != null) {
                    this.$panel.append(chart.getElem());
                    chart.init();
                    this.charts.push(chart);
                }
            }
            removeGraph(name) {
                let index = -1;
                for (let i = 0; i < this.charts.length; i++) {
                    let c = this.charts[i];
                    if (c.getOptions().name === name) {
                        index = i;
                        break;
                    }
                }
                if (index >= 0) {
                    let $elem = this.charts[index].getElem();
                    this.charts.splice(index, 1);
                    $elem.remove();
                }
            }
            save() {
                let charts = [];
                this.$panel.children().each((index, elem) => {
                    let name = $(elem).data("name");
                    for (let i = 0; i < this.charts.length; i++) {
                        let c = this.charts[i];
                        if (c.getOptions().name === name) {
                            charts.push({
                                name: c.getOptions().name,
                                width: c.getOptions().width,
                                height: c.getOptions().height,
                            });
                            break;
                        }
                    }
                });
                let args = {
                    name: this.opts.name,
                    key: this.opts.key || '',
                    charts: charts,
                };
                $ajax.post(`/system/chart/save_dashboard`, args).json((r) => {
                    if (r.success) {
                        Core.Notification.show("success", "Successfully saved.");
                    }
                    else {
                        Core.Notification.show("danger", r.message);
                    }
                });
            }
            getOptions() {
                return this.opts;
            }
            loadData() {
                if (this.charts.length == 0) {
                    return;
                }
                let args = {
                    charts: this.charts.map(c => c.getOptions().name).join(","),
                    period: this.opts.period,
                };
                if (this.opts.key) {
                    args.key = this.opts.key;
                }
                $ajax.get(`/system/chart/data`, args).json((d) => {
                    $.each(this.charts, (i, g) => {
                        if (d[g.getOptions().name]) {
                            g.setData(d[g.getOptions().name]);
                        }
                    });
                });
            }
        }
        Core.ChartDashboard = ChartDashboard;
        class ChartDialog {
            constructor(dashboard) {
                this.dashboard = dashboard;
                this.fb = new Core.FilterBox("#txt-query", this.filterCharts.bind(this));
                $("#btn-add").click(this.showAddDlg.bind(this));
                $("#btn-add-chart").click(this.addChart.bind(this));
                $("#btn-save").click(() => {
                    this.dashboard.save();
                });
            }
            showAddDlg() {
                let $panel = $("#nav-charts");
                $panel.find("label.panel-block").remove();
                $ajax.get(`/system/chart/query`, { dashboard: this.dashboard.getOptions().name }).json((charts) => {
                    for (let i = 0; i < charts.length; i++) {
                        let c = charts[i];
                        $panel.append(`<label class="panel-block">
          <input type="checkbox" value="${c.name}" data-index="${i}">${c.name}: ${c.title}
        </label>`);
                    }
                    this.charts = charts;
                    this.$charts = $panel.find("label.panel-block");
                });
                let dlg = new Core.Modal("#dlg-add-chart");
                dlg.show();
            }
            filterCharts(text) {
                if (!text) {
                    this.$charts.show();
                    return;
                }
                this.$charts.each((i, elem) => {
                    let $elem = $(elem);
                    let texts = [
                        this.charts[i].name.toLowerCase(),
                        this.charts[i].title.toLowerCase(),
                        this.charts[i].desc.toLowerCase(),
                    ];
                    for (let i = 0; i < texts.length; i++) {
                        let index = texts[i].indexOf(text);
                        if (index >= 0) {
                            $elem.show();
                            return;
                        }
                    }
                    $elem.hide();
                });
            }
            addChart() {
                this.$charts.each((i, e) => {
                    if ($(e).find(":checked").length > 0) {
                        let c = this.charts[i];
                        this.dashboard.addGraph(c);
                    }
                });
                Core.Modal.close();
            }
        }
    })(Core = Swirl.Core || (Swirl.Core = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var ChartDashboard = Swirl.Core.ChartDashboard;
    class IndexPage {
        constructor() {
            this.dashboard = new ChartDashboard("#div-charts", window.charts, { name: "home" });
            $("#cb-time").change(e => {
                this.dashboard.setPeriod($(e.target).val());
            });
            dragula([$('#div-charts').get(0)], {
                moves: function (el, container, handle) {
                    return $(handle).closest('a.drag').length > 0;
                }
            });
        }
    }
    Swirl.IndexPage = IndexPage;
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Metric;
    (function (Metric) {
        var Modal = Swirl.Core.Modal;
        var Dispatcher = Swirl.Core.Dispatcher;
        var FilterBox = Swirl.Core.FilterBox;
        class ListPage {
            constructor() {
                this.$charts = $("#div-charts").children();
                this.fb = new FilterBox("#txt-query", this.filterCharts.bind(this));
                $("#btn-import").click(this.importChart);
                Dispatcher.bind("#div-charts")
                    .on("export-chart", this.exportChart.bind(this))
                    .on("delete-chart", this.deleteChart.bind(this));
            }
            deleteChart(e) {
                let $container = $(e.target).closest("div.column");
                let name = $container.data("name");
                Modal.confirm(`Are you sure to delete chart: <strong>${name}</strong>?`, "Delete chart", (dlg, e) => {
                    $ajax.post(name + "/delete").trigger(e.target).json(r => {
                        $container.remove();
                        dlg.close();
                    });
                });
            }
            filterCharts(text) {
                if (!text) {
                    this.$charts.show();
                    return;
                }
                this.$charts.each((i, elem) => {
                    let $elem = $(elem), texts = [
                        $elem.data("name").toLowerCase(),
                        $elem.data("title").toLowerCase(),
                        $elem.data("desc").toLowerCase(),
                    ];
                    for (let i = 0; i < texts.length; i++) {
                        let index = texts[i].indexOf(text);
                        if (index >= 0) {
                            $elem.show();
                            return;
                        }
                    }
                    $elem.hide();
                });
            }
            exportChart(e) {
                let $container = $(e.target).closest("div.column");
                let name = $container.data("name");
                $ajax.get(name + "/detail").text(r => {
                    Modal.alert(`<textarea class="textarea" rows="8" readonly>${r}</textarea>`, "Export chart");
                });
            }
            importChart(e) {
                Modal.confirm(`<textarea class="textarea" rows="8"></textarea>`, "Import chart", (dlg, e) => {
                    try {
                        let chart = JSON.parse(dlg.find('textarea').val());
                        $ajax.post("new", chart).trigger(e.target).json(r => {
                            if (r.success) {
                                location.reload();
                            }
                            else {
                                dlg.error(r.message);
                            }
                        });
                    }
                    catch (e) {
                        dlg.error(e);
                    }
                });
            }
        }
        Metric.ListPage = ListPage;
    })(Metric = Swirl.Metric || (Swirl.Metric = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Config;
    (function (Config) {
        var Modal = Swirl.Core.Modal;
        var Table = Swirl.Core.ListTable;
        class ListPage {
            constructor() {
                this.table = new Table("#table-items");
                this.table.on("delete-config", this.deleteConfig.bind(this));
                $("#btn-delete").click(this.deleteConfigs.bind(this));
            }
            deleteConfig(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:eq(1)").text().trim();
                let id = $tr.find(":checkbox:first").val();
                Modal.confirm(`Are you sure to remove config: <strong>${name}</strong>?`, "Delete config", (dlg, e) => {
                    $ajax.post("delete", { ids: id }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            deleteConfigs(e) {
                let ids = this.table.selectedKeys();
                if (ids.length == 0) {
                    Modal.alert("Please select one or more items.");
                    return;
                }
                Modal.confirm(`Are you sure to remove ${ids.length} configs?`, "Delete configs", (dlg, e) => {
                    $ajax.post("delete", { ids: ids.join(",") }).trigger(e.target).encoder("form").json(r => {
                        this.table.selectedRows().remove();
                        dlg.close();
                    });
                });
            }
        }
        Config.ListPage = ListPage;
    })(Config = Swirl.Config || (Swirl.Config = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Config;
    (function (Config) {
        var OptionTable = Swirl.Core.OptionTable;
        class NewPage {
            constructor() {
                new OptionTable("#table-labels");
            }
        }
        Config.NewPage = NewPage;
    })(Config = Swirl.Config || (Swirl.Config = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Container;
    (function (Container) {
        var Modal = Swirl.Core.Modal;
        var Table = Swirl.Core.ListTable;
        class ListPage {
            constructor() {
                this.table = new Table("#table-items");
                this.table.on("delete-container", this.deleteContainer.bind(this));
                $("#btn-delete").click(this.deleteContainers.bind(this));
            }
            deleteContainer(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:eq(1)").text().trim();
                let id = $tr.find(":checkbox:first").val();
                Modal.confirm(`Are you sure to remove container: <strong>${name}</strong>?`, "Delete container", (dlg, e) => {
                    $ajax.post("delete", { ids: id }).trigger(e.target).encoder("form").json(() => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            deleteContainers() {
                let ids = this.table.selectedKeys();
                if (ids.length == 0) {
                    Modal.alert("Please select one or more items.");
                    return;
                }
                Modal.confirm(`Are you sure to remove ${ids.length} containers?`, "Delete containers", (dlg, e) => {
                    $ajax.post("delete", { ids: ids.join(",") }).trigger(e.target).encoder("form").json(() => {
                        this.table.selectedRows().remove();
                        dlg.close();
                    });
                });
            }
        }
        Container.ListPage = ListPage;
    })(Container = Swirl.Container || (Swirl.Container = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Image;
    (function (Image) {
        var Modal = Swirl.Core.Modal;
        var Table = Swirl.Core.ListTable;
        class ListPage {
            constructor() {
                this.table = new Table("#table-items");
                this.table.on("delete-image", this.deleteImage.bind(this));
                $("#btn-delete").click(this.deleteImages.bind(this));
            }
            deleteImage(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:eq(1)").text().trim();
                let id = $tr.find(":checkbox:first").val();
                Modal.confirm(`Are you sure to remove image: <strong>${name}</strong>?`, "Delete image", (dlg, e) => {
                    $ajax.post("delete", { ids: id }).trigger(e.target).encoder("form").json(() => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            deleteImages() {
                let ids = this.table.selectedKeys();
                if (ids.length == 0) {
                    Modal.alert("Please select one or more items.");
                    return;
                }
                Modal.confirm(`Are you sure to remove ${ids.length} images?`, "Delete images", (dlg, e) => {
                    $ajax.post("delete", { ids: ids.join(",") }).trigger(e.target).encoder("form").json(() => {
                        this.table.selectedRows().remove();
                        dlg.close();
                    });
                });
            }
        }
        Image.ListPage = ListPage;
    })(Image = Swirl.Image || (Swirl.Image = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Network;
    (function (Network) {
        var Modal = Swirl.Core.Modal;
        var Dispatcher = Swirl.Core.Dispatcher;
        class DetailPage {
            constructor() {
                let dispatcher = Dispatcher.bind("#table-containers");
                dispatcher.on("disconnect", this.disconnect.bind(this));
            }
            disconnect(e) {
                let $btn = $(e.target);
                let $tr = $btn.closest("tr");
                let id = $btn.val();
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to disconnect container: <strong>${name}</strong>?`, "Disconnect container", (dlg, e) => {
                    $ajax.post("disconnect", { container: id }).trigger($btn).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
        }
        Network.DetailPage = DetailPage;
    })(Network = Swirl.Network || (Swirl.Network = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Network;
    (function (Network) {
        var Modal = Swirl.Core.Modal;
        var Dispatcher = Swirl.Core.Dispatcher;
        class ListPage {
            constructor() {
                let dispatcher = Dispatcher.bind("#table-items");
                dispatcher.on("delete-network", this.deleteNetwork.bind(this));
            }
            deleteNetwork(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to remove network: <strong>${name}</strong>?`, "Delete network", (dlg, e) => {
                    $ajax.post("delete", { name: name }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
        }
        Network.ListPage = ListPage;
    })(Network = Swirl.Network || (Swirl.Network = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Network;
    (function (Network) {
        var OptionTable = Swirl.Core.OptionTable;
        class NewPage {
            constructor() {
                new OptionTable("#table-options");
                new OptionTable("#table-labels");
                $("#drivers :radio[name=driver]").change(e => {
                    $("#txt-custom-driver").prop("disabled", $(e.target).val() != "other");
                });
                $("#ipv6_enabled").change(e => {
                    let enabled = $(e.target).prop("checked");
                    $("#ipv6_subnet,#ipv6_gateway").prop("disabled", !enabled);
                });
            }
        }
        Network.NewPage = NewPage;
    })(Network = Swirl.Network || (Swirl.Network = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Node;
    (function (Node) {
        var OptionTable = Swirl.Core.OptionTable;
        class EditPage {
            constructor() {
                new OptionTable("#table-labels");
            }
        }
        Node.EditPage = EditPage;
    })(Node = Swirl.Node || (Swirl.Node = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Node;
    (function (Node) {
        var Modal = Swirl.Core.Modal;
        var Dispatcher = Swirl.Core.Dispatcher;
        class ListPage {
            constructor() {
                let dispatcher = Dispatcher.bind("#table-items");
                dispatcher.on("delete-node", this.deleteNode.bind(this));
            }
            deleteNode(e) {
                let $btn = $(e.target);
                let $tr = $btn.closest("tr");
                let id = $btn.val();
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to remove node: <strong>${name}</strong>?`, "Delete node", (dlg, e) => {
                    $ajax.post("delete", { id: id }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
        }
        Node.ListPage = ListPage;
    })(Node = Swirl.Node || (Swirl.Node = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Perm;
    (function (Perm) {
        var Dispatcher = Swirl.Core.Dispatcher;
        var Modal = Swirl.Core.Modal;
        class EditPage {
            constructor() {
                $("#txt-query").keydown(this.searchUser);
                $("#btn-add-user").click(this.addUser);
                Dispatcher.bind("#div-users").on("delete-user", this.deleteUser.bind(this));
            }
            deleteUser(e) {
                $(e.target).closest("div.control").remove();
            }
            searchUser(e) {
                if (e.keyCode == 13) {
                    let query = $.trim($(e.target).val());
                    if (query.length == 0) {
                        return;
                    }
                    $ajax.post("/system/user/search", { query: query }).encoder("form").json((users) => {
                        let $panel = $("#nav-users");
                        $panel.find("label.panel-block").remove();
                        for (let user of users) {
                            $panel.append(`<label class="panel-block">
          <input type="checkbox" value="${user.id}" data-name="${user.name}"> ${user.name}
        </label>`);
                        }
                    });
                }
            }
            addUser() {
                let users = {};
                $("#div-users").find("input").each((i, e) => {
                    users[$(e).val()] = true;
                });
                let $panel = $("#nav-users");
                $panel.find("input:checked").each((i, e) => {
                    let $el = $(e);
                    if (users[$el.val()]) {
                        return;
                    }
                    $("#div-users").append(`<div class="control">
            <div class="tags has-addons">
              <span class="tag is-info">${$el.data("name")}</span>
              <a class="tag is-delete" data-action="delete-user"></a>
              <input name="users[]" value="${$el.val()}" type="hidden">
            </div>`);
                });
                Modal.close();
                $panel.find("label.panel-block").remove();
            }
        }
        Perm.EditPage = EditPage;
    })(Perm = Swirl.Perm || (Swirl.Perm = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Registry;
    (function (Registry) {
        var Modal = Swirl.Core.Modal;
        var Dispatcher = Swirl.Core.Dispatcher;
        class ListPage {
            constructor() {
                let dispatcher = Dispatcher.bind("#table-items");
                dispatcher.on("delete-registry", this.deleteRegistry.bind(this));
                dispatcher.on("edit-registry", this.editRegistry.bind(this));
            }
            deleteRegistry(e) {
                let $tr = $(e.target).closest("tr");
                let id = $tr.data("id");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to remove registry: <strong>${name}</strong>?`, "Delete registry", (dlg, e) => {
                    $ajax.post("delete", { id: id }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            editRegistry(e) {
                let $tr = $(e.target).closest("tr");
                let dlg = new Modal("#dlg-edit");
                dlg.find("input[name=id]").val($tr.data("id"));
                dlg.find("input[name=name]").val($tr.find("td:first").text().trim());
                dlg.find("input[name=url]").val($tr.find("td:eq(1)").text().trim());
                dlg.find("input[name=username]").val($tr.find("td:eq(2)").text().trim());
                dlg.show();
            }
        }
        Registry.ListPage = ListPage;
    })(Registry = Swirl.Registry || (Swirl.Registry = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Role;
    (function (Role) {
        var Dispatcher = Swirl.Core.Dispatcher;
        class EditPage {
            constructor() {
                $("#table-perms").find("tr").each((i, elem) => {
                    let $tr = $(elem);
                    let $cbs = $tr.find("td :checkbox");
                    $tr.find("th>:checkbox").prop("checked", $cbs.length == $cbs.filter(":checked").length);
                });
                Dispatcher.bind("#table-perms", "change")
                    .on("check-row", this.checkRow.bind(this))
                    .on("check", this.check.bind(this));
            }
            checkRow(e) {
                let $cb = $(e.target);
                let checked = $cb.prop("checked");
                $cb.closest("th").next("td").find(":checkbox").prop("checked", checked);
            }
            check(e) {
                let $cb = $(e.target);
                let $cbs = $cb.closest("td").find(":checkbox");
                let checked = $cbs.length == $cbs.filter(":checked").length;
                $cb.closest("td").prev("th").find(":checkbox").prop("checked", checked);
            }
        }
        Role.EditPage = EditPage;
    })(Role = Swirl.Role || (Swirl.Role = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Role;
    (function (Role) {
        var Modal = Swirl.Core.Modal;
        var Dispatcher = Swirl.Core.Dispatcher;
        class ListPage {
            constructor() {
                Dispatcher.bind("#table-items")
                    .on("delete-role", this.deleteUser.bind(this));
            }
            deleteUser(e) {
                let $tr = $(e.target).closest("tr");
                let id = $tr.data("id");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to remove role: <strong>${name}</strong>?`, "Delete role", (dlg, e) => {
                    $ajax.post("delete", { id: id }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
        }
        Role.ListPage = ListPage;
    })(Role = Swirl.Role || (Swirl.Role = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Role;
    (function (Role) {
        var Dispatcher = Swirl.Core.Dispatcher;
        class NewPage {
            constructor() {
                Dispatcher.bind("#table-perms", "change")
                    .on("check-row", this.checkRow.bind(this))
                    .on("check", this.check.bind(this));
            }
            checkRow(e) {
                let $cb = $(e.target);
                let checked = $cb.prop("checked");
                $cb.closest("th").next("td").find(":checkbox").prop("checked", checked);
            }
            check(e) {
                let $cb = $(e.target);
                let $cbs = $cb.closest("td").find(":checkbox");
                let checked = $cbs.length == $cbs.filter(":checked").length;
                $cb.closest("td").prev("th").find(":checkbox").prop("checked", checked);
            }
        }
        Role.NewPage = NewPage;
    })(Role = Swirl.Role || (Swirl.Role = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Secret;
    (function (Secret) {
        var Modal = Swirl.Core.Modal;
        var Table = Swirl.Core.ListTable;
        class ListPage {
            constructor() {
                this.table = new Table("#table-items");
                this.table.on("delete-secret", this.deleteSecret.bind(this));
                $("#btn-delete").click(this.deleteSecrets.bind(this));
            }
            deleteSecret(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:eq(1)").text().trim();
                let id = $tr.find(":checkbox:first").val();
                Modal.confirm(`Are you sure to remove secret: <strong>${name}</strong>?`, "Delete secret", (dlg, e) => {
                    $ajax.post("delete", { ids: id }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            deleteSecrets(e) {
                let ids = this.table.selectedKeys();
                if (ids.length == 0) {
                    Modal.alert("Please select one or more items.");
                    return;
                }
                Modal.confirm(`Are you sure to remove ${ids.length} secrets?`, "Delete secrets", (dlg, e) => {
                    $ajax.post("delete", { ids: ids.join(",") }).trigger(e.target).encoder("form").json(r => {
                        this.table.selectedRows().remove();
                        dlg.close();
                    });
                });
            }
        }
        Secret.ListPage = ListPage;
    })(Secret = Swirl.Secret || (Swirl.Secret = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Secret;
    (function (Secret) {
        var OptionTable = Swirl.Core.OptionTable;
        class NewPage {
            constructor() {
                new OptionTable("#table-labels");
            }
        }
        Secret.NewPage = NewPage;
    })(Secret = Swirl.Secret || (Swirl.Secret = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Service;
    (function (Service) {
        var Validator = Swirl.Core.Validator;
        var OptionTable = Swirl.Core.OptionTable;
        var EditTable = Swirl.Core.EditTable;
        var Modal = Swirl.Core.Modal;
        var Table = Swirl.Core.Table;
        class ServiceModeRule {
            constructor($model) {
                this.$mode = $model;
            }
            validate($form, $input, arg) {
                if (this.$mode.val() == "global") {
                    return { ok: true };
                }
                const regex = /^(0|[1-9]\d*)$/;
                return { ok: regex.test($.trim($input.val())) };
            }
        }
        class MountTable extends EditTable {
            render() {
                return `<tr>
              <td>
                <div class="select is-small">
                  <select name="mounts[${this.index}].type">
                    <option value="bind">Bind</option>
                    <option value="volume">Volume</option>
                    <option value="tmpfs">TempFS</option>
                  </select>
                </div>
              </td>
              <td>
                <input name="mounts[${this.index}].src" class="input is-small" placeholder="path in host">
              </td>
              <td><input name="mounts[${this.index}].dst" class="input is-small" placeholder="path in container"></td>
              <td>
                <div class="select is-small">
                  <select name="mounts[${this.index}].read_only" data-type="bool">
                    <option value="false">No</option>
                    <option value="true">Yes</option>
                  </select>
                </div>
              </td>
              <td>
                <div class="select is-small">
                  <select name="mounts[${this.index}].propagation">
                    <option value="">--Select--</option>
                    <option value="rprivate">rprivate</option>
                    <option value="private">private</option>
                    <option value="rshared">rshared</option>
                    <option value="shared">shared</option>
                    <option value="rslave">rslave</option>
                    <option value="slave">slave</option>
                  </select>
                </div>
              </td>
              <td>
                <a class="button is-small is-outlined is-danger" data-action="delete-mount">
                  <span class="icon is-small">
                    <i class="far fa-trash-alt"></i>
                  </span>
                </a>
              </td>
            </tr>`;
            }
        }
        class PortTable extends EditTable {
            render() {
                return `<tr>
                <td><input name="endpoint.ports[${this.index}].published_port" class="input is-small" placeholder="port in host" data-type="integer"></td>
                <td>
                  <input name="endpoint.ports[${this.index}].target_port" class="input is-small" placeholder="port in container" data-type="integer">
                </td>
                <td>
                  <div class="select is-small">
                    <select name="endpoint.ports[${this.index}].protocol">
                      <option value="false">TCP</option>
                      <option value="true">UDP</option>
                    </select>
                  </div>
                </td>
                <td>
                  <div class="select is-small">
                    <select name="endpoint.ports[${this.index}].publish_mode">
                      <option value="ingress">ingress</option>
                      <option value="host">host</option>
                    </select>
                  </div>
                </td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-endpoint-port">
                    <span class="icon is-small">
                      <i class="far fa-trash-alt"></i>
                    </span>
                  </a>
                </td>
              </tr>`;
            }
        }
        class ConstraintTable extends EditTable {
            render() {
                return `<tr>
                <td>
                  <input name="placement.constraints[${this.index}].name" class="input is-small" placeholder="e.g. node.role/node.hostname/node.id/node.labels.*/engine.labels.*/...">
                </td>
                <td>
                  <div class="select is-small">
                    <select name="placement.constraints[${this.index}].op">
                      <option value="==">==</option>
                      <option value="!=">!=</option>
                    </select>
                  </div>
                </td>
                <td>
                  <input name="placement.constraints[${this.index}].value" class="input is-small" placeholder="e.g. manager">
                </td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-constraint">
                    <span class="icon is-small">
                      <i class="far fa-trash-alt"></i>
                    </span>
                  </a>
                </td>
              </tr>`;
            }
        }
        class PreferenceTable extends EditTable {
            render() {
                return `<tr>
                <td>
                  <input name="placement.preferences[${this.index}].spread" class="input is-small" placeholder="e.g. engine.labels.az">
                </td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-preference">
                    <span class="icon is-small">
                      <i class="far fa-trash-alt"></i>
                    </span>
                  </a>
                </td>
              </tr>`;
            }
        }
        class ConfigTable extends Table {
            constructor(elem) {
                super(elem);
                this.name = this.$table.data("name");
                this.$body = this.$table.find("tbody");
                this.index = this.$body.find("tr").length;
                super.on("add-" + this.name, this.showAddDialog.bind(this)).on("delete-" + this.name, ConfigTable.deleteRow);
            }
            addRow(id, name) {
                let field = `${this.name}s[${this.index}]`;
                this.$body.append(`<tr>
                <td>${name}<input name="${field}.id" value="${id}" type="hidden"><input name="${field}.name" value="${name}" type="hidden"></td>
                <td><input name="${field}.file_name" value="${name}" class="input is-small"></td>
                <td><input name="${field}.uid" value="0" class="input is-small"></td>
                <td><input name="${field}.gid" value="0" class="input is-small"></td>
                <td><input name="${field}.mode" value="444" class="input is-small" data-type="integer"></td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-${this.name}">
                    <span class="icon is-small">
                      <i class="far fa-trash-alt"></i>
                    </span>
                  </a>
                </td>
              </tr>`);
                this.index++;
            }
            showAddDialog(e) {
                let dlg = new Modal("#dlg-add-" + this.name);
                dlg.find(":checked").prop("checked", false);
                dlg.error();
                dlg.show();
            }
            static deleteRow(e) {
                $(e.target).closest("tr").remove();
            }
        }
        class EditPage {
            constructor() {
                this.$mode = $("#cb-mode");
                this.$replicas = $("#txt-replicas");
                new OptionTable("#table-envs");
                new OptionTable("#table-slabels");
                new OptionTable("#table-clabels");
                new OptionTable("#table-log_driver-options");
                new PortTable("#table-endpoint-ports");
                new MountTable("#table-mounts");
                new ConstraintTable("#table-constraints");
                new PreferenceTable("#table-preferences");
                this.secret = new ConfigTable("#table-secrets");
                this.config = new ConfigTable("#table-configs");
                Validator.register("service-mode", new ServiceModeRule(this.$mode), "Please input a valid integer.");
                this.$mode.change(e => this.$replicas.toggle(this.$mode.val() != "global"));
                $("#btn-add-secret").click(() => EditPage.addConfig(this.secret));
                $("#btn-add-config").click(() => EditPage.addConfig(this.config));
            }
            static addConfig(t) {
                let dlg = Modal.current();
                let $cbs = dlg.find(":checked");
                if ($cbs.length == 0) {
                    dlg.error(`Please select the ${t.name} files.`);
                }
                else {
                    dlg.close();
                    $cbs.each((i, cb) => {
                        let $cb = $(cb);
                        t.addRow($cb.val(), $cb.data("name"));
                    });
                }
            }
        }
        Service.EditPage = EditPage;
        class NewPage extends EditPage {
            constructor() {
                super();
                this.$registryUrl = $("#a-registry-url");
                this.$registry = $("#cb-registry");
                this.$registry.change(e => {
                    let url = this.$registry.find("option:selected").data("url") || "";
                    this.$registryUrl.text(url + "/").toggle(url != "");
                });
            }
        }
        Service.NewPage = NewPage;
    })(Service = Swirl.Service || (Swirl.Service = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Service;
    (function (Service) {
        var Modal = Swirl.Core.Modal;
        var Table = Swirl.Core.ListTable;
        class ListPage {
            constructor() {
                this.table = new Table("#table-items");
                this.table.on("delete-service", this.deleteService.bind(this))
                    .on("scale-service", this.scaleService.bind(this))
                    .on("rollback-service", this.rollbackService.bind(this))
                    .on("restart-service", this.restartService.bind(this));
            }
            deleteService(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:eq(0)").text().trim();
                Modal.confirm(`Are you sure to remove service: <strong>${name}</strong>?`, "Delete service", (dlg, e) => {
                    $ajax.post(`${name}/delete`, { names: name }).trigger(e.target).encoder("form").json(() => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            scaleService(e) {
                let $target = $(e.target), $tr = $target.closest("tr");
                let data = {
                    name: $tr.find("td:eq(0)").text().trim(),
                    count: $target.data("replicas"),
                };
                Modal.confirm(`<input name="count" value="${data.count}" class="input" placeholder="Replicas">`, "Scale service", dlg => {
                    data.count = dlg.find("input[name=count]").val();
                    $ajax.post(`${data.name}/scale`, data).encoder("form").json(() => {
                        location.reload();
                    });
                });
            }
            rollbackService(e) {
                let $tr = $(e.target).closest("tr"), name = $tr.find("td:eq(0)").text().trim();
                Modal.confirm(`Are you sure to rollback service: <strong>${name}</strong>?`, "Rollback service", dlg => {
                    $ajax.post(`${name}/rollback`, { name: name }).encoder("form").json(() => {
                        dlg.close();
                    });
                });
            }
            restartService(e) {
                let $tr = $(e.target).closest("tr"), name = $tr.find("td:eq(0)").text().trim();
                Modal.confirm(`Are you sure to restart service: <strong>${name}</strong>?`, "Restart service", dlg => {
                    $ajax.post(`${name}/restart`, { name: name }).encoder("form").json(() => {
                        dlg.close();
                    });
                });
            }
        }
        Service.ListPage = ListPage;
    })(Service = Swirl.Service || (Swirl.Service = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Service;
    (function (Service) {
        var ChartDashboard = Swirl.Core.ChartDashboard;
        class StatsPage {
            constructor() {
                let $cb_time = $("#cb-time");
                if ($cb_time.length == 0) {
                    return;
                }
                this.dashboard = new ChartDashboard("#div-charts", window.charts, {
                    name: "service",
                    key: $("#h2-service-name").text()
                });
                dragula([$('#div-charts').get(0)]);
                $cb_time.change(e => {
                    this.dashboard.setPeriod($(e.target).val());
                });
                $("#cb-refresh").change(e => {
                    if ($(e.target).prop("checked")) {
                        this.dashboard.refresh();
                    }
                    else {
                        this.dashboard.stop();
                    }
                });
            }
        }
        Service.StatsPage = StatsPage;
    })(Service = Swirl.Service || (Swirl.Service = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Service;
    (function (Service) {
        var Template;
        (function (Template) {
            var Modal = Swirl.Core.Modal;
            var Dispatcher = Swirl.Core.Dispatcher;
            class ListPage {
                constructor() {
                    let dispatcher = Dispatcher.bind("#table-items");
                    dispatcher.on("delete-template", this.deleteTemplate.bind(this));
                }
                deleteTemplate(e) {
                    let $tr = $(e.target).closest("tr");
                    let id = $tr.data("id");
                    let name = $tr.find("td:first").text();
                    Modal.confirm(`Are you sure to remove template: <strong>${name}</strong>?`, "Delete template", (dlg, e) => {
                        $ajax.post("delete", { id: id }).trigger(e.target).encoder("form").json(() => {
                            $tr.remove();
                            dlg.close();
                        });
                    });
                }
            }
            Template.ListPage = ListPage;
        })(Template = Service.Template || (Service.Template = {}));
    })(Service = Swirl.Service || (Swirl.Service = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Setting;
    (function (Setting) {
        class IndexPage {
            constructor() {
                $("#ldap-enabled").change(e => {
                    let enabled = $(e.target).prop("checked");
                    $("#fs-ldap").find("input:not(:checkbox)").prop("readonly", !enabled);
                });
                $("#ldap-auth-simple,#ldap-auth-bind").click(e => {
                    if ($(e.target).val() == "0") {
                        $("#div-auth-simple").show();
                        $("#div-auth-bind").hide();
                    }
                    else {
                        $("#div-auth-simple").hide();
                        $("#div-auth-bind").show();
                    }
                });
            }
        }
        Setting.IndexPage = IndexPage;
    })(Setting = Swirl.Setting || (Swirl.Setting = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Stack;
    (function (Stack) {
        var Validator = Swirl.Core.Validator;
        var Notification = Swirl.Core.Notification;
        class ContentRequiredRule {
            validate($form, $input, arg) {
                let el = $input[0];
                if ($("#type-" + arg).prop("checked")) {
                    return { ok: el.checkValidity ? el.checkValidity() : true, error: el.validationMessage };
                }
                return { ok: true };
            }
        }
        class EditPage {
            constructor() {
                Validator.register("content", new ContentRequiredRule(), "");
                this.editor = CodeMirror.fromTextArea($("#txt-content")[0], { lineNumbers: true });
                $("#file-content").change(e => {
                    let file = e.target;
                    if (file.files.length > 0) {
                        $('#filename').text(file.files[0].name);
                    }
                });
                $("#type-input,#type-upload").click(e => {
                    let type = $(e.target).val();
                    $("#div-input").toggle(type == "input");
                    $("#div-upload").toggle(type == "upload");
                });
                $("#btn-submit").click(this.submit.bind(this));
            }
            submit(e) {
                this.editor.save();
                let results = Validator.bind("#div-form").validate();
                if (results.length > 0) {
                    return;
                }
                let data = new FormData();
                data.append('name', $("#name").val());
                if ($("#type-input").prop("checked")) {
                    data.append('content', $('#txt-content').val());
                }
                else {
                    let file = $('#file-content')[0];
                    data.append('content', file.files[0]);
                }
                let url = $(e.target).data("url") || "";
                $ajax.post(url, data).encoder("none").trigger(e.target).json((r) => {
                    if (r.success) {
                        location.href = "/stack/";
                    }
                    else {
                        Notification.show("danger", `FAILED: ${r.message}`);
                    }
                });
            }
        }
        Stack.EditPage = EditPage;
    })(Stack = Swirl.Stack || (Swirl.Stack = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Stack;
    (function (Stack) {
        var Modal = Swirl.Core.Modal;
        var Dispatcher = Swirl.Core.Dispatcher;
        class ListPage {
            constructor() {
                let dispatcher = Dispatcher.bind("#table-items");
                dispatcher.on("deploy-stack", this.deployStack.bind(this));
                dispatcher.on("shutdown-stack", this.shutdownStack.bind(this));
                dispatcher.on("delete-stack", this.deleteStack.bind(this));
            }
            deployStack(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to deploy stack: <strong>${name}</strong>?`, "Deploy stack", (dlg, e) => {
                    $ajax.post(`${name}/deploy`).trigger(e.target).json(r => {
                        location.reload();
                    });
                });
            }
            shutdownStack(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to shutdown stack: <strong>${name}</strong>?`, "Shutdown stack", (dlg, e) => {
                    $ajax.post(`${name}/shutdown`).trigger(e.target).json(r => {
                        location.reload();
                    });
                });
            }
            deleteStack(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to remove archive: <strong>${name}</strong>?`, "Delete stack", (dlg, e) => {
                    $ajax.post(`${name}/delete`).trigger(e.target).json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
        }
        Stack.ListPage = ListPage;
    })(Stack = Swirl.Stack || (Swirl.Stack = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Task;
    (function (Task) {
        class LogsPage {
            constructor() {
                this.refreshInterval = 3000;
                this.$line = $("#txt-line");
                this.$timestamps = $("#cb-timestamps");
                this.$refresh = $("#cb-refresh");
                this.$stdout = $("#txt-stdout");
                this.$stderr = $("#txt-stderr");
                this.$refresh.change(e => {
                    let elem = (e.target);
                    if (elem.checked) {
                        this.refreshData();
                    }
                    else if (this.timer > 0) {
                        window.clearTimeout(this.timer);
                        this.timer = 0;
                    }
                });
                this.refreshData();
            }
            refreshData() {
                let args = {
                    line: this.$line.val(),
                    timestamps: this.$timestamps.prop("checked"),
                };
                $ajax.get('fetch_logs', args).json((r) => {
                    this.$stdout.val(r.stdout);
                    this.$stderr.val(r.stderr);
                    this.$stdout.get(0).scrollTop = this.$stdout.get(0).scrollHeight;
                    this.$stderr.get(0).scrollTop = this.$stderr.get(0).scrollHeight;
                });
                if (this.$refresh.prop("checked")) {
                    this.timer = setTimeout(this.refreshData.bind(this), this.refreshInterval);
                }
            }
        }
        Task.LogsPage = LogsPage;
    })(Task = Swirl.Task || (Swirl.Task = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var User;
    (function (User) {
        var Modal = Swirl.Core.Modal;
        var Dispatcher = Swirl.Core.Dispatcher;
        class ListPage {
            constructor() {
                Dispatcher.bind("#table-items")
                    .on("delete-user", this.deleteUser.bind(this))
                    .on("block-user", this.blockUser.bind(this))
                    .on("unblock-user", this.unblockUser.bind(this));
            }
            deleteUser(e) {
                let $tr = $(e.target).closest("tr");
                let id = $tr.data("id");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to remove user: <strong>${name}</strong>?`, "Delete user", (dlg, e) => {
                    $ajax.post("delete", { id: id }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            blockUser(e) {
                let $tr = $(e.target).closest("tr");
                let id = $tr.data("id");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to block user: <strong>${name}</strong>?`, "Block user", (dlg, e) => {
                    $ajax.post("block", { id: id }).trigger(e.target).encoder("form").json(r => {
                        location.reload();
                    });
                });
            }
            unblockUser(e) {
                let $tr = $(e.target).closest("tr");
                let id = $tr.data("id");
                let name = $tr.find("td:first").text().trim();
                Modal.confirm(`Are you sure to unblock user: <strong>${name}</strong>?`, "Unblock user", (dlg, e) => {
                    $ajax.post("unblock", { id: id }).trigger(e.target).encoder("form").json(r => {
                        location.reload();
                    });
                });
            }
        }
        User.ListPage = ListPage;
    })(User = Swirl.User || (Swirl.User = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Volume;
    (function (Volume) {
        var Modal = Swirl.Core.Modal;
        var Table = Swirl.Core.ListTable;
        class ListPage {
            constructor() {
                this.table = new Table("#table-items");
                this.table.on("delete-volume", this.deleteVolume.bind(this));
                $("#btn-delete").click(this.deleteVolumes.bind(this));
                $("#btn-prune").click(this.pruneVolumes.bind(this));
            }
            deleteVolume(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:eq(1)").text().trim();
                Modal.confirm(`Are you sure to remove volume: <strong>${name}</strong>?`, "Delete volume", (dlg, e) => {
                    $ajax.post("delete", { names: name }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            deleteVolumes(e) {
                let names = this.table.selectedKeys();
                if (names.length == 0) {
                    Modal.alert("Please select one or more items.");
                    return;
                }
                Modal.confirm(`Are you sure to remove ${names.length} volumes?`, "Delete volumes", (dlg, e) => {
                    $ajax.post("delete", { names: names.join(",") }).trigger(e.target).encoder("form").json(r => {
                        this.table.selectedRows().remove();
                        dlg.close();
                    });
                });
            }
            pruneVolumes(e) {
                Modal.confirm(`Are you sure to remove all unused volumes?`, "Prune volumes", (dlg, e) => {
                    $ajax.post("prune").trigger(e.target).json(r => {
                        location.reload();
                    });
                });
            }
        }
        Volume.ListPage = ListPage;
    })(Volume = Swirl.Volume || (Swirl.Volume = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Volume;
    (function (Volume) {
        var OptionTable = Swirl.Core.OptionTable;
        class NewPage {
            constructor() {
                new OptionTable("#table-options");
                new OptionTable("#table-labels");
                $("#drivers").find(":radio[name=driver]").change(e => {
                    $("#txt-custom-driver").prop("disabled", $(e.target).val() != "other");
                });
            }
        }
        Volume.NewPage = NewPage;
    })(Volume = Swirl.Volume || (Swirl.Volume = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Container;
    (function (Container) {
        class ExecPage {
            constructor() {
                this.$cmd = $("#txt-cmd");
                this.$connect = $("#btn-connect");
                this.$disconnect = $("#btn-disconnect");
                this.$connect.click(this.connect.bind(this));
                this.$disconnect.click(this.disconnect.bind(this));
                Terminal.applyAddon(fit);
            }
            connect(e) {
                this.$connect.hide();
                this.$disconnect.show();
                let url = location.host + location.pathname.substring(0, location.pathname.lastIndexOf("/")) + "/connect?cmd=" + encodeURIComponent(this.$cmd.val());
                let ws = new WebSocket("ws://" + url);
                ws.onopen = () => {
                    this.term = new Terminal();
                    this.term.on('data', (data) => {
                        if (ws.readyState == WebSocket.OPEN) {
                            ws.send(data);
                        }
                    });
                    this.term.open(document.getElementById('terminal-container'));
                    this.term.focus();
                    let width = Math.floor(($('#terminal-container').width() - 20) / 8.39);
                    let height = 30;
                    this.term.resize(width, height);
                    this.term.setOption('cursorBlink', true);
                    this.term.fit();
                    window.onresize = () => {
                        this.term.fit();
                    };
                    ws.onmessage = (e) => {
                        this.term.write(e.data);
                    };
                    ws.onerror = function (error) {
                        console.log("error: " + error);
                    };
                    ws.onclose = () => {
                        console.log("close");
                    };
                };
                this.ws = ws;
            }
            disconnect(e) {
                if (this.ws && this.ws.readyState != WebSocket.CLOSED) {
                    this.ws.close();
                }
                if (this.term) {
                    this.term.destroy();
                    this.term = null;
                }
                this.$connect.show();
                this.$disconnect.hide();
            }
        }
        Container.ExecPage = ExecPage;
    })(Container = Swirl.Container || (Swirl.Container = {}));
})(Swirl || (Swirl = {}));
//# sourceMappingURL=swirl.js.map