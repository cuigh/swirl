///<reference path="../core/core.ts" />
namespace Swirl.Container {
    export class ExecPage {
        private $cmd: JQuery;
        private $connect: JQuery;
        private $disconnect: JQuery;
        private ws: WebSocket;
        private term: Terminal;

        constructor() {
            this.$cmd = $("#txt-cmd");
            this.$connect = $("#btn-connect");
            this.$disconnect = $("#btn-disconnect");

            this.$connect.click(this.connect.bind(this));
            this.$disconnect.click(this.disconnect.bind(this));

            Terminal.applyAddon(fit);
        }

        private connect(e: JQueryEventObject) {
            this.$connect.hide();
            this.$disconnect.show();

            let url = location.host + location.pathname.substring(0, location.pathname.lastIndexOf("/")) + "/connect?cmd=" + encodeURIComponent(this.$cmd.val());
            let protocol = (location.protocol === "https:") ? "wss://" : "ws://";
            let ws = new WebSocket(protocol + url);
            ws.onopen = () => {
                this.term = new Terminal();
                this.term.on('data', (data: any) => {
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

        private disconnect(e: JQueryEventObject) {
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
}