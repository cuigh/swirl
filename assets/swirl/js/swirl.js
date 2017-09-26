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
                var doubleByteChars = value.match(/[^\x00-\xff]/ig);
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
                let regex = new RegExp(arg, 'i');
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
                <span class="icon is-small"><i class="fa fa-trash"></i></span>
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
                const regex = /^[1-9]\d*$/;
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
                    <i class="fa fa-trash"></i>
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
                      <option value="rprivate">ingress</option>
                      <option value="private">host</option>
                    </select>
                  </div>
                </td>
                <td>
                  <a class="button is-small is-outlined is-danger" data-action="delete-endpoint-port">
                    <span class="icon is-small">
                      <i class="fa fa-trash"></i>
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
                      <i class="fa fa-trash"></i>
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
                      <i class="fa fa-trash"></i>
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
                      <i class="fa fa-remove"></i>
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
                Validator.register("service-mode", new ServiceModeRule(this.$mode), "Please input a positive integer.");
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
                this.table.on("delete-service", this.deleteService.bind(this)).on("scale-service", this.scaleService.bind(this));
                $("#btn-delete").click(this.deleteServices.bind(this));
            }
            deleteService(e) {
                let $tr = $(e.target).closest("tr");
                let name = $tr.find("td:eq(1)").text().trim();
                Modal.confirm(`Are you sure to remove service: <strong>${name}</strong>?`, "Delete service", (dlg, e) => {
                    $ajax.post("delete", { names: name }).trigger(e.target).encoder("form").json(r => {
                        $tr.remove();
                        dlg.close();
                    });
                });
            }
            deleteServices(e) {
                let names = this.table.selectedKeys();
                if (names.length == 0) {
                    Modal.alert("Please select one or more items.");
                    return;
                }
                Modal.confirm(`Are you sure to remove ${names.length} services?`, "Delete services", (dlg, e) => {
                    $ajax.post("delete", { names: names.join(",") }).trigger(e.target).encoder("form").json(r => {
                        this.table.selectedRows().remove();
                        dlg.close();
                    });
                });
            }
            scaleService(e) {
                let $btn = $(e.target);
                let $tr = $btn.closest("tr");
                let data = {
                    name: $tr.find("td:eq(1)").text().trim(),
                    count: $btn.data("replicas"),
                };
                Modal.confirm(`<input name="count" value="${data.count}" class="input" placeholder="Replicas">`, "Scale service", (dlg, e) => {
                    data.count = dlg.find("input[name=count]").val();
                    $ajax.post("scale", data).trigger($btn).encoder("form").json(r => {
                        location.reload();
                    });
                });
            }
        }
        Service.ListPage = ListPage;
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
            }
        }
        Setting.IndexPage = IndexPage;
    })(Setting = Swirl.Setting || (Swirl.Setting = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Stack;
    (function (Stack) {
        var Archive;
        (function (Archive) {
            var Modal = Swirl.Core.Modal;
            var Dispatcher = Swirl.Core.Dispatcher;
            class ListPage {
                constructor() {
                    let dispatcher = Dispatcher.bind("#table-items");
                    dispatcher.on("deploy-archive", this.deployArchive.bind(this));
                    dispatcher.on("delete-archive", this.deleteArchive.bind(this));
                }
                deployArchive(e) {
                    let $tr = $(e.target).closest("tr");
                    let id = $tr.data("id");
                    let name = $tr.find("td:first").text().trim();
                    Modal.confirm(`Are you sure to deploy archive: <strong>${name}</strong>?`, "Deploy archive", (dlg, e) => {
                        $ajax.post("deploy", { id: id }).trigger(e.target).encoder("form").json(r => {
                            dlg.close();
                        });
                    });
                }
                deleteArchive(e) {
                    let $tr = $(e.target).closest("tr");
                    let id = $tr.data("id");
                    let name = $tr.find("td:first").text().trim();
                    Modal.confirm(`Are you sure to remove archive: <strong>${name}</strong>?`, "Delete archive", (dlg, e) => {
                        $ajax.post("delete", { id: id }).trigger(e.target).encoder("form").json(r => {
                            $tr.remove();
                            dlg.close();
                        });
                    });
                }
            }
            Archive.ListPage = ListPage;
        })(Archive = Stack.Archive || (Stack.Archive = {}));
    })(Stack = Swirl.Stack || (Swirl.Stack = {}));
})(Swirl || (Swirl = {}));
var Swirl;
(function (Swirl) {
    var Stack;
    (function (Stack) {
        var Task;
        (function (Task) {
            var Modal = Swirl.Core.Modal;
            var Dispatcher = Swirl.Core.Dispatcher;
            class ListPage {
                constructor() {
                    let dispatcher = Dispatcher.bind("#table-items");
                    dispatcher.on("delete-stack", this.deleteStack.bind(this));
                }
                deleteStack(e) {
                    let $tr = $(e.target).closest("tr");
                    let name = $tr.find("td:first").text().trim();
                    Modal.confirm(`Are you sure to remove stack: <strong>${name}</strong>?`, "Delete stack", (dlg, e) => {
                        $ajax.post("delete", { name: name }).trigger(e.target).encoder("form").json(r => {
                            $tr.remove();
                            dlg.close();
                        });
                    });
                }
            }
            Task.ListPage = ListPage;
        })(Task = Stack.Task || (Stack.Task = {}));
    })(Stack = Swirl.Stack || (Swirl.Stack = {}));
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
//# sourceMappingURL=swirl.js.map