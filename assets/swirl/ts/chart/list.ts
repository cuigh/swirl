///<reference path="../core/core.ts" />
namespace Swirl.Metric {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    class FilterBox {
        private $elem: JQuery;
        private timer: number;

        constructor(elem: string | Element | JQuery, callback: (text: string) => void, timeout?: number) {
            this.$elem = $(elem);
            this.$elem.keyup(() => {
                if (this.timer > 0) {
                    clearTimeout(this.timer);
                }
                let text: string = this.$elem.val().toLowerCase();
                this.timer = setTimeout(() => callback(text), timeout || 500);
            });
        }
    }

    export class ListPage {
        private fb: FilterBox;
        private $charts: JQuery;

        constructor() {
            this.$charts = $("#div-charts").children();
            this.fb = new FilterBox("#txt-query", this.filterCharts.bind(this));

            // bind events
            Dispatcher.bind("#div-charts").on("delete-chart", this.deleteChart.bind(this));
        }

        private deleteChart(e: JQueryEventObject) {
            let $container = $(e.target).closest("div.column");
            let name = $container.data("name");
            Modal.confirm(`Are you sure to delete chart: <strong>${name}</strong>?`, "Delete chart", (dlg, e) => {
                $ajax.post(name + "/delete").trigger(e.target).json<AjaxResult>(r => {
                    $container.remove();
                    dlg.close();
                })
            });
        }

        private filterCharts(text: string) {
            if (!text) {
                this.$charts.show();
                return;
            }

            this.$charts.each((i, elem) => {
                let $elem = $(elem),
                    texts: string[] = [
                        $elem.data("name").toLowerCase(),
                        $elem.data("title").toLowerCase(),
                        $elem.data("desc").toLowerCase(),
                    ];
                for (let i = 0; i<texts.length; i++) {
                    let index = texts[i].indexOf(text);
                    if (index >= 0) {
                        $(elem).show();
                        return;
                    }
                }
                $(elem).hide();
            })
        }
    }
}