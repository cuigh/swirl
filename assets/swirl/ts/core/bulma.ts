namespace Swirl.Core {
    export class Modal {
        private static initialized: boolean;
        private static active: Modal;
        private $modal: JQuery;
        private deletable: boolean;

        constructor(modal: string | Element | JQuery) {
            this.$modal = $(modal);
            this.find(".modal-background, .modal-close, .modal-card-head .delete, .modal-card-foot .dismiss").click(e => this.close());
        }

        static initialize() {
            if (!Modal.initialized) {
                $('.modal-trigger').click(function () {
                    let target = $(this).data('target');
                    let dlg = new Modal("#" + target)
                    dlg.show();
                });
                $(document).on('keyup', function (e) {
                    // ESC Key
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

        static current(): Modal {
            return Modal.active;
        }

        static alert(content: string, title?: string, callback?: (dlg: Modal, e: JQueryEventObject) => void): Modal {
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
            callback = callback || function (dlg: Modal, e: JQueryEventObject) { dlg.close(); };
            $elem.find(".modal-card-foot>button:first").click(e => callback(dlg, e));
            dlg.deletable = true;
            dlg.show();
            return dlg;
        }

        static confirm(content: string, title?: string, callback?: (dlg: Modal, e: JQueryEventObject) => void): Modal {
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
            } else {
                this.$modal.removeClass('is-active');
            }
            Modal.active = null;
        }

        error(msg?: string) {
            let $error = this.find(".modal-card-error");
            if (msg) {
                if (!$error.length) {
                    $error = $('<section class="modal-card-error"></section>').insertAfter(this.find(".modal-card-body"));
                }
                $error.html(msg).show();
            } else {
                $error.hide();
            }
        }

        find(selector: string | Element | JQuery): JQuery {
            if (typeof selector === "string") {
                return this.$modal.find(selector);
            } else if (selector instanceof Element) {
                return this.$modal.find(selector);
            }
            return this.$modal.find(selector);
        }
    }

    export type NotificationType = "primary" | "info" | "success" | "warning" | "danger";

    export interface NotificationOptions {
        type: NotificationType;
        time: number; // display time(in seconds), 0 is always hide
        message: string;
    }

    export class Notification {
        private $elem: JQuery;
        private options: NotificationOptions;

        constructor(options: NotificationOptions) {
            this.options = $.extend({}, {
                type: "primary",
                time: 5,
            }, options);
        }

        static show(type: NotificationType, msg: string, time?: number): Notification {
            let n = new Notification({
                type: type,
                message: msg,
                time: time,
            });
            n.show();
            return n;
        }

        private show() {
            this.$elem = $(`<div class="notification is-${this.options.type} has-text-centered is-marginless is-radiusless is-top" style="display:none">${this.options.message}</div>`).appendTo("body").fadeIn(250);
            if (this.options.time > 0) {
                setTimeout(() => this.hide(), this.options.time * 1000);
            }
        }

        hide() {
            this.$elem.fadeTo("slow", 0.01, () => this.$elem.remove());
        }
    }

    export class Tab {
        private static initialized: boolean;
        private $tab: JQuery;
        private $content: JQuery;
        private active: number;

        constructor(tab: string | Element | JQuery, content: string | Element | JQuery) {
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

        select(index: number) {
            if (this.active != index) {
                this.$tab.find(`li.is-active`).removeClass("is-active");
                this.$tab.find(`li:eq(${index})`).addClass("is-active");
                this.$content.children(":visible").hide();
                $(this.$content.children()[index]).show();
                this.active = index;
            }
        }
    }
}
