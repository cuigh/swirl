///<reference path="../core/core.ts" />
namespace Swirl.Metric {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;
    import FilterBox = Swirl.Core.FilterBox;

    export class ListPage {
        private fb: FilterBox;
        private $charts: JQuery;

        constructor() {
            this.$charts = $("#div-charts").children();
            this.fb = new FilterBox("#txt-query", this.filterCharts.bind(this));

            // bind events
            $("#btn-import").click(this.importChart);
            Dispatcher.bind("#div-charts")
                .on("export-chart", this.exportChart.bind(this))
                .on("delete-chart", this.deleteChart.bind(this));
        }

        private deleteChart(e: JQueryEventObject) {
            let $container = $(e.target).closest("div.column");
            let name = $container.data("name");
            Modal.confirm(`Are you sure to delete chart: <strong>${name}</strong>?`, "Delete chart", (dlg, e) => {
                $ajax.post(name + "/delete").trigger(e.target).json<AjaxResult>(r => {
                    $container.remove();
                    dlg.close();
                });
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
                for (let i = 0; i < texts.length; i++) {
                    let index = texts[i].indexOf(text);
                    if (index >= 0) {
                        $elem.show();
                        return;
                    }
                }
                $elem.hide();
            })
        }

        private exportChart(e: JQueryEventObject) {
            let $container = $(e.target).closest("div.column");
            let name = $container.data("name");
            $ajax.get(name + "/detail").text(r => {
                Modal.alert(`<textarea class="textarea" rows="8" readonly>${r}</textarea>`, "Export chart");
            });
        }

        private importChart(e: JQueryEventObject) {
            Modal.confirm(`<textarea class="textarea" rows="8"></textarea>`, "Import chart", (dlg, e) => {
                try {
                    let chart = JSON.parse(dlg.find('textarea').val());
                    $ajax.post("new", chart).trigger(e.target).json<AjaxResult>(r => {
                        if (r.success) {
                            location.reload();
                        } else {
                            dlg.error(r.message);
                        }
                    });
                } catch (e) {
                    dlg.error(e);
                }
            });
        }
    }
}