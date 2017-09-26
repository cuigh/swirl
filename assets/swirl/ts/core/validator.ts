/*!
 * Swirl Validator Library v1.0.0
 * Copyright 2017 cuigh. All rights reserved.
 *
 * @author cuigh(noname@live.com)
 */
namespace Swirl.Core {
    type InputElement = HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement | HTMLButtonElement;

    /**
     * 输入控件验证结果
     *
     * @interface ValidationResult
     */
    export interface ValidationResult {
        input: JQuery;
        errors: string[];
    }

    // interface ValidationRule {
    //     ($input: JQuery): boolean;
    // }

    export interface ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string};
    }

    export interface ValidationMarker {
        setError($input: JQuery, errors: string[]): void;
        clearError($input: JQuery): void;
        reset($input: JQuery): void;
    }

    /**
     * HTML5 表单元素原生验证器
     *
     * @class NativeRule
     * @implements {ValidationRule}
     */
    class NativeRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            let el = <InputElement>$input[0];
            return {ok: el.checkValidity ? el.checkValidity() : true};
        }
    }

    /**
     * 必填字段验证器
     *
     * @class RequiredRule
     * @implements {ValidationRule}
     */
    class RequiredRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            return {ok: $.trim($input.val()).length > 0};
        }
    }

    /**
     * 必选字段验证器(用于 radio 和 checkbox), 示例: checked, checked(2), checked(1~2)
     *
     * @class CheckedRule
     * @implements {ValidationRule}
     */
    class CheckedRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            let count = parseInt(arg);
            let siblings = $form.find(`:input:checked[name='${$input.attr("name")}']`);
            return {ok: siblings.length >= count};
        }
    }

    /**
     * 电子邮箱验证器
     *
     * @class EmailValidator
     * @implements {ValidationRule}
     */
    class EmailRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            const regex = /^((([a-z]|\d|[!#\$%&'\*\+\-\/=\?\^_`{\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+(\.([a-z]|\d|[!#\$%&'\*\+\-\/=\?\^_`{\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+)*)|((\x22)((((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(([\x01-\x08\x0b\x0c\x0e-\x1f\x7f]|\x21|[\x23-\x5b]|[\x5d-\x7e]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(\\([\x01-\x09\x0b\x0c\x0d-\x7f]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))))*(((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(\x22)))@((([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.)+(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))$/i;
            let value = $.trim($input.val());
            return {ok: !value || regex.test(value)};
        }
    }

    /**
     * HTTP/FTP 地址验证器
     *
     * @class UrlValidator
     * @implements {ValidationRule}
     */
    class UrlRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            const regex = /^(https?|ftp):\/\/(((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:)*@)?(((\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5]))|((([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.)+(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.?)(:\d*)?)(\/((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)+(\/(([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)*)*)?)?(\?((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)|[\uE000-\uF8FF]|\/|\?)*)?(\#((([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(%[\da-f]{2})|[!\$&'\(\)\*\+,;=]|:|@)|\/|\?)*)?$/i;
            let value = $.trim($input.val());
            return {ok: !value || regex.test(value)};
        }
    }

    /**
     * IPV4 地址验证器
     *
     * @class IPValidator
     * @implements {ValidationRule}
     */
    class IPRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            const regex = /^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[0-9]{1,2})$/i;
            let value = $.trim($input.val());
            return {ok: !value || regex.test(value)};
        }
    }

    /**
     * 字段匹配验证器(如密码)
     *
     * @class MatchValidator
     * @implements {ValidationRule}
     */
    class MatchValidator implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            return {ok: $input.val() == $('#' + arg).val()};
        }
    }

    /**
     * 字符串长度验证器
     *
     * @class LengthValidator
     * @implements {ValidationRule}
     */
    class LengthRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            let r: {ok: boolean, error?: string} = {ok: true};
            if (arg) {
                let len = this.getLength($.trim($input.val()));
                let args = arg.split('~');
                if (args.length == 1) {
                    if ($.isNumeric(args[0])) {
                        r.ok = len >= parseInt(args[0]);
                    }
                } else {
                    if ($.isNumeric(args[0]) && $.isNumeric(args[1])) {
                        r.ok = len >= parseInt(args[0]) && len <= parseInt(args[1])
                    }
                }
            }
            return r;
        }

        protected getLength(value: string): number {
            return value.length;
        }
    }

    /**
     * 字符串宽度验证器(中文字符宽度为2)
     *
     * @class WidthValidator
     * @extends {LengthRule}
     */
    class WidthRule extends LengthRule {
        protected getLength(value: string): number {
            var doubleByteChars = value.match(/[^\x00-\xff]/ig);
            return value.length + (doubleByteChars == null ? 0 : doubleByteChars.length);
        }
    }

    /**
     * 整数验证器
     *
     * @class IntegerValidator
     * @implements {ValidationRule}
     */
    class IntegerRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            const regex = /^\d*$/;
            return {ok: regex.test($.trim($input.val()))};
        }
    }

    /**
     * 正则表达式验证器
     *
     * @class RegexValidator
     * @implements {ValidationRule}
     */
    class RegexRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            let regex = new RegExp(arg, 'i');
            let value = $.trim($input.val());
            return {ok: !value || regex.test(value)};
        }
    }

    /**
     * 服务器端验证器
     *
     * @class RemoteRule
     * @implements {ValidationRule}
     */
    class RemoteRule implements ValidationRule {
        validate($form: JQuery, $input: JQuery, arg?: string): {ok: boolean, error?: string} {
            if (!arg) {
                throw new Error("服务器验证地址未设置");
            }

            let value = $.trim($input.val());
            let r: {ok: boolean, error?: string} = {ok: false};
            $ajax.post(arg, {value: value}).encoder("form").async(false).json<{error: string}>(result => {
                r.ok = !result.error;
                r.error = result.error;
            });
            return r;
        }
    }

    export interface ValidatorOptions {
    }

    export class Validator {
        private static selector = ':input[data-v-rule]:not(:submit,:button,:reset,:image,:disabled)';
        // error marker
        private static marker: ValidationMarker;
        // error message
        private static messages: { [index: string]: string } = {
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
        private static rules: { [index: string]: ValidationRule } = {
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
        private form: JQuery;
        private options: Object;

        /**
         * Creates an instance of Validator.
         *
         * @param {(string | HTMLElement | JQuery)} elem the parent element which contains all form inputs
         * @param {*} [options] the validation options
         *
         * @memberOf Validator
         */
        private constructor(elem: string | HTMLElement | JQuery, options?: ValidatorOptions) {
            this.form = $(elem);
            this.options = options;

            // disable default validation of HTML5, and bind submit event
            if (this.form.is("form")) {
                this.form.attr("novalidate", "true");
                // this.form.submit(e => {
                //     let results = this.validate();
                //     if (results != null && results.length > 0) {
                //         e.preventDefault();
                //     }
                // });
            }

            // realtime validate events
            this.form.on("click", ':radio[data-v-rule],:checkbox[data-v-rule]', this.checkValue.bind(this));
            this.form.on("change", 'select[data-v-rule],input[type="file"][data-v-rule]', this.checkValue.bind(this));
            this.form.on("blur", ':input[data-v-rule]:not(select,:radio,:checkbox,:file)', this.checkValue.bind(this));
        }

        private checkValue(e: JQueryEventObject) {
            let $input = $(e.target);
            let result = this.validateInput($input);
            Validator.mark($input, result);
        }

        /**
         * 创建验证器并绑定到表单
         *
         * @static
         * @param {(string | HTMLElement | JQuery)} elem 验证表单或其它容器元素
         * @param {ValidatorOptions} [options]
         * @returns {Validator} 选项
         *
         * @memberOf Validator
         */
        static bind(elem: string | HTMLElement | JQuery, options?: ValidatorOptions): Validator {
            let v = $(elem).data("validator");
            if (!v) {
                v = new Validator(elem, options);
                $(elem).data("validator", v);
            }
            return v;
        }

        /**
         * 验证表单
         *
         * @returns {ValidationResult[]}
         */
        validate(): ValidationResult[] {
            let results: ValidationResult[] = [];
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

        /**
         * 清除验证标识
         */
        reset(): void {
            this.form.find(Validator.selector).each((i, el) => {
                let $input = $(el);
                Validator.marker.reset($input);
            });
        }

        /**
         * 注册验证器
         *
         * @static
         * @param {string} name 验证器名称
         * @param {ValidationRule} rule 验证方法
         * @param {string} msg 验证消息
         *
         * @memberOf Validator
         */
        static register(name: string, rule: ValidationRule, msg: string) {
            this.rules[name] = rule;
            this.messages[name] = msg;
        }

        /**
         * set error message
         */
        static setMessage(name: string, msg: string) {
            this.messages[name] = msg;
        }

        /**
         * set error marker
         */
        static setMarker(marker: ValidationMarker) {
            this.marker = marker;
        }

        private validateInput($input: JQuery): ValidationResult {
            let errors: string[] = [];
            let rules: string[]= ($input.data('v-rule') || 'native').split(';');
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

        private static mark($input: JQuery, result: ValidationResult) {
            if (Validator.marker != null) {
                if (result) {
                    Validator.marker.setError($input, result.errors);
                }
                else {
                    Validator.marker.clearError($input);
                }
            }
        }

        private static getMessge($input: JQuery, rule: string) {
            // $input[0].validationMessage
            // if (!success) $input[0].setCustomValidity("错误信息");
            if (rule == 'native') return (<InputElement>$input[0]).validationMessage;
            else {
                let msg = $input.data('v-msg-' + rule);
                if (!msg) msg = this.messages[rule];
                return msg;
            }
        }
    }
}
