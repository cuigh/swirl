namespace Swirl.Core {
        /**
     * Dispatcher
     */
    export class Dispatcher {
        private attr: string;

        private events: { [index: string]: (e: JQueryEventObject) => any } = {};

        constructor(attr?: string) {
            this.attr = attr || "action";
        }

        /**
         * 创建一个 Dispatcher 并绑定事件到页面元素上
         *
         * @param elem
         * @param event
         * @returns {Dispatcher}
         */
        static bind(elem: string | JQuery | Element | Document, event: string = "click"): Dispatcher {
            return new Dispatcher().bind(elem, event);
        }

        /**
         * 注册动作事件
         *
         * @param action
         * @param handler
         * @returns {Mtime.Util.Dispatcher}
         */
        on(action: string, handler: (e: JQueryEventObject) => any): Dispatcher {
            this.events[action] = handler;
            return this;
        }

        /**
         * 移除动作事件
         *
         * @param action
         * @returns {Mtime.Util.Dispatcher}
         */
        off(action: string): Dispatcher {
            delete this.events[action];
            return this;
        }

        /**
         * 绑定事件到页面元素上
         *
         * @param elem
         * @param event
         * @returns {Mtime.Util.Dispatcher}
         */
        bind(elem: string | JQuery | Element | Document, event: string = "click"): Dispatcher {
            $(elem).on(event, this.handle.bind(this));
            return this;
        }

        private handle(e: JQueryEventObject): any {
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
}