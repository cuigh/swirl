///<reference path="../core/core.ts" />
namespace Swirl.Task {
    export class LogsPage {
        private $stdout: JQuery;
        private $stderr: JQuery;
        private $line: JQuery;
        private $timestamps: JQuery;
        private $refresh: JQuery;
        private timer: number;
        private refreshInterval = 3000;

        constructor() {
            this.$line = $("#txt-line");
            this.$timestamps = $("#cb-timestamps");
            this.$refresh = $("#cb-refresh");
            this.$stdout = $("#txt-stdout");
            this.$stderr = $("#txt-stderr");

            this.$refresh.change(e => {
                let elem = <HTMLInputElement>(e.target);
                if (elem.checked) {
                    this.refreshData();
                } else if (this.timer > 0) {
                    window.clearTimeout(this.timer);
                    this.timer = 0;
                }
            });

            this.refreshData();
            // let ws = new WebSocket("ws://" + location.host + location.pathname + "_ws");
            // ws.onopen = e => console.log("open");
            // ws.onclose = e => console.log("close");
            // ws.onmessage = e => console.log("message: " + e.data);
        }

        private refreshData() {
            let args: any = {
                line: this.$line.val(),
                timestamps: this.$timestamps.prop("checked"),
            };
            $ajax.get('fetch_logs', args).json((r: any) => {
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
}