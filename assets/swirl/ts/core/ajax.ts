/*!
 * Swirl Ajax Library v1.0.0
 * Copyright 2017 cuigh. All rights reserved.
 *
 * @author cuigh(noname@live.com)
 */
///<reference path="bulma.ts"/>
namespace Swirl.Core {
    export type AjaxErrorHandler = (xhr: JQueryXHR, textStatus: string, error: string) => void;

    /**
     * AjaxResult
     */
    export class AjaxResult {
        success: boolean;
        code?: number;
        message?: string;
        data?: any;
        url?: string;
    }

    /**
     * AjaxOptions
     */
    export class AjaxOptions {
        private static defaultOptions: AjaxOptions = new AjaxOptions();
        /**
         * request url
         */
        url: string;
        /**
         * request method, GET/POST...
         */
        method: AjaxMethod;
        /**
         * request data
         */
        data?: Object;
        /**
         * request timeout time(ms)
         *
         * @type {number}
         */
        timeout?: number = 30000;
        /**
         * send request by asynchronous
         *
         * @type {boolean}
         */
        async?: boolean = true;
        /**
         * response data type
         */
        dataType?: "text" | "html" | "json" | "jsonp" | "xml" | "script" | string;
        /**
         * AJAX trigger element for indicator
         *
         * @type {(Element | JQuery)}
         */
        trigger?: Element | JQuery;
        /**
         * data encoder for POST request
         */
        encoder: "none" | "form" | "json" = "json";
        /**
         * previous filter
         */
        preHandler: (options: AjaxOptions) => void;
        /**
         * post filter
         */
        postHandler: (options: AjaxOptions) => void;
        /**
         * error handler
         */
        errorHandler: AjaxErrorHandler;

        /**
         * get default options
         *
         * @returns {AjaxOptions}
         */
        static getDefaultOptions(): AjaxOptions {
            return AjaxOptions.defaultOptions;
        }

        /**
         * set default options
         *
         * @param options
         */
        static setDefaultOptions(options: AjaxOptions) {
            if (options) {
                AjaxOptions.defaultOptions = options;
            }
        }
    }

    /**
     * AJAX request method
     */
    export enum AjaxMethod {
        GET,
        POST,
        PUT,
        DELETE,
        HEAD,
        TRACE,
        OPTIONS,
        CONNECT,
        PATCH
    }

    /**
     * AJAX Request
     */
    export class AjaxRequest {
        static preHandler: (options: AjaxOptions) => void = options => {
            options.trigger && $(options.trigger).prop("disabled", true);
        };
        static postHandler: (options: AjaxOptions) => void = options => {
            options.trigger && $(options.trigger).prop("disabled", false);
        };
        static errorHandler: (xhr: JQueryXHR, textStatus: string, error: string) => void = (xhr, status, err) => {
            let msg: string;
            if (xhr.responseJSON) {
                // auxo web framework return: {code: 0, message: "xxx"}
                let err = xhr.responseJSON;
                msg = err.message;
                if (err.code) {
                    msg += `(${err.code})`
                }                
            } else if (xhr.status >= 400) {
                msg = xhr.responseText || err || status;
            } else {
                return
            }
            Notification.show("danger", `AJAX: ${msg}`)            
        };
        protected options: AjaxOptions;

        constructor(url: string, method: AjaxMethod, data?: any) {
            this.options = $.extend({
                url: url,
                method: method,
                data: data,
                preHandler: AjaxRequest.preHandler,
                postHandler: AjaxRequest.postHandler,
                errorHandler: AjaxRequest.errorHandler,
            }, AjaxOptions.getDefaultOptions());
        }

        /**
         * set pre handler
         *
         * @param handler
         * @return {AjaxRequest}
         */
        preHandler(handler: (options: AjaxOptions) => void): this {
            this.options.preHandler = handler;
            return this;
        }

        /**
         * set post handler
         *
         * @param handler
         * @return {AjaxRequest}
         */
        postHandler(handler: (options: AjaxOptions) => void): this {
            this.options.postHandler = handler;
            return this;
        }

        /**
         * set error handler
         *
         * @param handler
         * @return {AjaxRequest}
         */
        errorHandler(handler: AjaxErrorHandler): this {
            this.options.errorHandler = handler;
            return this;
        }

        /**
         * set request timeout
         *
         * @param timeout
         * @returns {AjaxRequest}
         */
        timeout(timeout: number): this {
            this.options.timeout = timeout;
            return this;
        }

        /**
         * set async option
         *
         * @param async
         * @returns {AjaxRequest}
         */
        async(async: boolean): this {
            this.options.async = async;
            return this;
        }

        /**
         * set trigger element
         *
         * @param {(Element | JQuery)} elem
         * @returns {AjaxRequest}
         */
        trigger(elem: Element | JQuery): this {
            this.options.trigger = elem;
            return this;
        }

        /**
         * get response as JSON
         *
         * @template T JSON data type
         * @param {(r: T) => void} [callback] callback function
         */
        json<T>(callback?: (r: T) => void): void | Promise<T> {
            return this.result<T>("json", callback);
        }

        /**
         * get response as text
         *
         * @param {(r: string) => void} [callback] callback function
         */
        text(callback?: (r: string) => void): void | Promise<string> {
            return this.result<string>("text", callback);
        }

        /**
         * get response as HTML
         *
         * @param {(r: string) => void} [callback] callback function
         */
        html(callback?: (r: string) => void): void | Promise<string> {
            return this.result<string>("html", callback);
        }

        protected result<T>(dataType: string, callback?: (r: T) => void): void | Promise<T> {
            this.options.dataType = dataType;
            this.options.preHandler && this.options.preHandler(this.options);
            let settings = this.buildSettings();
            if (callback) {
                $.ajax(settings).done(callback).always(() => {
                    this.options.postHandler && this.options.postHandler(this.options);
                }).fail((xhr: JQueryXHR, textStatus: string, error: string) => {
                    this.options.errorHandler && this.options.errorHandler(xhr, textStatus, error);
                });
            } else {
                return new Promise<T>((resolve, _) => {
                    $.ajax(settings).done((r: T) => {
                        resolve(r);
                    }).always(() => {
                        AjaxRequest.postHandler && AjaxRequest.postHandler(this.options);
                    }).fail((xhr: JQueryXHR, textStatus: string, error: string) => {
                        AjaxRequest.errorHandler && AjaxRequest.errorHandler(xhr, textStatus, error);
                    });
                });
            }
        }

        protected buildSettings(): JQueryAjaxSettings {
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

    /**
     * AJAX GET Request
     */
    export class AjaxGetRequest extends AjaxRequest {
        /**
         * get JSON response by jsonp
         *
         * @template T JSON data type
         * @param {(r: T) => void} [callback] callback function
         */
        jsonp<T>(callback?: (r: T) => void): void | Promise<T> {
            return this.result<T>("jsonp", callback);
        }
    }

    /**
     * AJAX POST Request
     */
    export class AjaxPostRequest extends AjaxRequest {
        /**
         * set encoder
         *
         * @param encoder
         * @returns {AjaxPostRequest}
         */
        encoder(encoder: "none" | "form" | "json"): this {
            this.options.encoder = encoder;
            return this;
        }

        protected buildSettings(): JQueryAjaxSettings {
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

    /**
     * AJAX static entry class
     */
    export class Ajax {
        private constructor() {
        }

        /**
         * Send GET request
         *
         * @static
         * @param {string} url request url
         * @param {Object} [args] request data
         * @returns {Ajax} Ajax request instance
         */
        static get(url: string, args?: Object): AjaxGetRequest {
            return new AjaxGetRequest(url, AjaxMethod.GET, args);
        }

        /**
         * Send POST request
         *
         * @static
         * @param {string} url request url
         * @param {(string | Object)} [data] request data
         * @returns {Ajax} Ajax request instance
         */
        static post(url: string, data?: string | Object): AjaxPostRequest {
            return new AjaxPostRequest(url, AjaxMethod.POST, data);
        }
    }

    /**
     * AJAX interface(仅用于实现 $ajax 快捷对象)
     */
    export interface AjaxStatic {
        /**
         * Send GET request
         *
         * @static
         * @param {string} url request url
         * @param {Object} [args] request data
         * @returns {Ajax} Ajax request instance
         */
        get(url: string, args?: Object): AjaxGetRequest;
        /**
         * Send POST request
         *
         * @static
         * @param {string} url request url
         * @param {(string | Object)} [data] request data
         * @returns {Ajax} Ajax request instance
         */
        post(url: string, data?: string | Object): AjaxPostRequest;
    }
}

let $ajax: Swirl.Core.AjaxStatic = Swirl.Core.Ajax;
