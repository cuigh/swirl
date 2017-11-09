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
         * Create a Dispatcher instance
         *
         * @param elem
         * @param event
         * @returns {Dispatcher}
         */
        static bind(elem: string | JQuery | Element | Document, event: string = "click"): Dispatcher {
            return new Dispatcher().bind(elem, event);
        }

        /**
         * Register event
         *
         * @param action
         * @param handler
         * @returns {Swirl.Core.Dispatcher}
         */
        on(action: string, handler: (e: JQueryEventObject) => any): Dispatcher {
            this.events[action] = handler;
            return this;
        }

        /**
         * Unregister event
         *
         * @param action
         * @returns {Swirl.Core.Dispatcher}
         */
        off(action: string): Dispatcher {
            delete this.events[action];
            return this;
        }

        /**
         * Bind events to element
         *
         * @param elem
         * @param event
         * @returns {Swirl.Core.Dispatcher}
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