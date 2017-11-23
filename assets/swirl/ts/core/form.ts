/*!
 * Swirl Form Library v1.0.0
 * Copyright 2017 cuigh. All rights reserved.
 * see also: https://github.com/A1rPun/transForm.js
 *
 * @author cuigh(noname@live.com)
 */
///<reference path="validator.ts"/>
namespace Swirl.Core {
    /**
     * Form options
     */
    export class FormOptions {
        delimiter: string = ".";
        skipDisabled: boolean = true;
        skipReadOnly: boolean = false;
        skipEmpty: boolean = false;
        useIdOnEmptyName: boolean = true;
        triggerChange: boolean = false;
    }

    interface Entry {
        name: string;
        value?: any;
    }

    /**
     * Form
     */
    export class Form {
        private form: JQuery;
        private options: FormOptions;
        private validator: Validator;

        /**
         * Creates an instance of Form.
         *
         * @param {(string | Element | JQuery)} elem Form html element
         * @param {FormOptions} options Form options
         *
         * @memberOf Form
         */
        constructor(elem: string | Element | JQuery, options?: FormOptions) {
            this.form = $(elem);
            this.options = $.extend(new FormOptions(), options);
            this.validator = Validator.bind(this.form);
        }

        /**
         * Reset form
         *
         * @memberOf Form
         */
        reset() {
            (<HTMLFormElement>this.form.get(0)).reset();
            if (this.validator) {
                this.validator.reset();
            }
        }

        /**
         * Clear form data
         *
         * @memberOf Form
         */
        clear() {
            let inputs = this.getFields();
            inputs.each((i, input) => {
                this.clearInput(input);
            });
            if (this.validator) {
                this.validator.reset();
            }
        }

        /**
         * Submit form by AJAX
         *
         * @param {string} url submit url
         * @returns {Swirl.Core.AjaxPostRequest}
         *
         * @memberOf Form
         */
        submit(url?: string): AjaxPostRequest {
            let data = this.serialize();
            return Ajax.post(url || this.form.attr("action"), data);
        }

        /**
         * Validate form
         *
         * @returns {boolean}
         */
        validate(): boolean {
            if (!this.validator) {
                return true;
            }
            return this.validator.validate().length == 0;
        }

        /**
         * Serialize form data to JSON
         *
         * @param {Function} [nodeCallback] custom callback for parsing input value
         * @returns {Object}
         *
         * @memberOf Form
         */
        serialize(nodeCallback?: Function): Object {
            let result = {},
                inputs = this.getFields();

            for (let i = 0, l = inputs.length; i < l; i++) {
                let input: any = inputs[i],
                    key = input.name || this.options.useIdOnEmptyName && input.id;

                if (!key) continue;

                let entry: Entry = null;
                if (nodeCallback) entry = nodeCallback(input);
                if (!entry) entry = this.getEntryFromInput(input, key);

                if (typeof entry.value === 'undefined' || entry.value === null
                    || (this.options.skipEmpty && (!entry.value || (this.isArray(entry.value) && !entry.value.length))))
                    continue;
                this.saveEntryToResult(result, entry, input, this.options.delimiter);
            }
            return result;
        }

        /**
         * Fill data to form
         *
         * @param {Object} data JSON data
         * @param {Function} [nodeCallback] custom callback for processing input value
         *
         * @memberOf Form
         */
        deserialize(data: Object, nodeCallback?: Function) {
            let inputs = this.getFields();
            for (let i = 0, l = inputs.length; i < l; i++) {
                let input: any = inputs[i],
                    key = input.name || this.options.useIdOnEmptyName && input.id,
                    value = this.getFieldValue(key, data);

                if (typeof value === 'undefined' || value === null) {
                    this.clearInput(input);
                    continue;
                }

                let mutated = nodeCallback && nodeCallback(input, value);
                if (!mutated) this.setValueToInput(input, value);
            }
        }

        /**
         * automatic initialize Form components using data attributes.
         */
        static automate() {
            $('form[data-form]').each(function (i, elem) {
                let $form = $(elem);
                let form = new Form($form);
                let type = $form.data("form");

                if (type == "form") {
                    $form.submit(e => { return form.validate() });
                    return;
                }

                // ajax-json | ajax-form
                $form.submit(function () {
                    if (!form.validate()) {
                        return false;
                    }

                    let request = form.submit($form.attr("action")).trigger($form.find('button[type="submit"]'));
                    if (type == "ajax-form") {
                        request.encoder("form");
                    }
                    request.json((r: AjaxResult) => {
                        if (r.success) {
                            let url = r.url || $form.data("url");
                            if (url) {
                                if (url === "-") {
                                    location.reload();
                                } else {
                                    location.href = url;
                                }
                            } else {
                                let msg = r.message || $form.data("message");
                                Notification.show("info", `SUCCESS: ${msg}`, 3);
                            }
                        } else {
                            let msg = r.message;
                            if (r.code) {
                                msg += `({r.code})`
                            }
                            Notification.show("danger", `FAILED: ${msg}`);
                        }
                    });
                    return false;
                });
            });
        }

        private getEntryFromInput(input: any, key: string): Entry {
            let nodeType = input.type && input.type.toLowerCase(),
                entry: Entry = { name: key },
                dataType = $(input).data("type");

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
                        if (input.options[i].selected) entry.value.push(input.options[i].value);
                    break;
                case 'file':
                    //Only interested in the filename (Chrome adds C:\fakepath\ for security anyway)
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

        private saveEntryToResult(parent: { [index: string]: any }, entry: Entry, input: HTMLElement, delimiter: string) {
            //not not accept empty values in array collections
            if (/\[]$/.test(entry.name) && !entry.value) return;

            let parts = this.parseString(entry.name, delimiter);
            for (let i = 0, l = parts.length; i < l; i++) {
                let part = parts[i];
                //if last
                if (i === l - 1) {
                    parent[part] = entry.value;
                } else {
                    let index = parts[i + 1];
                    if (!index || $.isNumeric(index)) {
                        if (!this.isArray(parent[part]))
                            parent[part] = [];
                        //if second last
                        if (i === l - 2) {
                            parent[part].push(entry.value);
                        } else {
                            if (!this.isObject(parent[part][index]))
                                parent[part][index] = {};
                            parent = parent[part][index];
                        }
                        i++;
                    } else {
                        if (!this.isObject(parent[part]))
                            parent[part] = {};
                        parent = parent[part];
                    }
                }
            }
        }

        private clearInput(input: any) {
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
                    if (input.checked) input.checked = false;
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

        private parseString(str: string, delimiter: string): string[] {
            let result: string[] = [],
                split = str.split(delimiter),
                len = split.length;
            for (let i = 0; i < len; i++) {
                let s = split[i].split('['),
                    l = s.length;
                for (let j = 0; j < l; j++) {
                    let key = s[j];
                    if (!key) {
                        //if the first one is empty, continue
                        if (j === 0) continue;
                        //if the undefined key is not the last part of the string, throw error
                        if (j !== l - 1)
                            throw new Error(`Undefined key is not the last part of the name > ${str}`);
                    }
                    //strip "]" if its there
                    if (key && key[key.length - 1] === ']')
                        key = key.slice(0, -1);
                    result.push(key);
                }
            }
            return result;
        }

        private getFields(): JQuery {
            let inputs: JQuery = this.form.find("input,select,textarea").filter(':not([data-form-ignore="true"])');
            if (this.options.skipDisabled) inputs = inputs.filter(":not([disabled])");
            if (this.options.skipReadOnly) inputs = inputs.filter(":not([readonly])");
            return inputs;
        }

        private getFieldValue(key: string, ref: any): any {
            if (!key || !ref) return;

            let parts = this.parseString(key, this.options.delimiter);
            for (let i = 0, l = parts.length; i < l; i++) {
                let part = ref[parts[i]];

                if (typeof part === 'undefined' || part === null) return;

                //if last
                if (i === l - 1) {
                    return part;
                } else {
                    let index = parts[i + 1];
                    if (index === '') {
                        return part;
                    } else if ($.isNumeric(index)) {
                        //if second last
                        if (i === l - 2)
                            return part[index];
                        else
                            ref = part[index];
                        i++;
                    } else {
                        ref = part;
                    }
                }
            }
        }

        private setValueToInput(input: any, value: any) {
            let nodeType = input.type && input.type.toLowerCase();
            switch (nodeType) {
                case 'radio':
                    if (value == input.value) input.checked = true;
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

        /*** Helper functions ***/

        private contains(array: any[], value: any): boolean {
            for (let item of array) {
                if (item == value) return true;
            }
            return false;
        }

        private isObject(obj: any) {
            return typeof obj === 'object';
        }

        private isArray(arr: any) {
            return Array.isArray(arr);
        }
    }
}