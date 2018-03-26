///<reference path="core/core.ts" />
///<reference path="core/graph.ts" />
namespace Swirl {
    import Modal = Swirl.Core.Modal;
    import GraphPanel = Swirl.Core.GraphPanel;
    import FilterBox = Swirl.Core.FilterBox;

    export class IndexPage {
        private panel: GraphPanel;
        private fb: FilterBox;
        private charts: any;
        private $charts: JQuery;

        constructor() {
            this.fb = new FilterBox("#txt-query", this.filterCharts.bind(this));
            this.panel = new GraphPanel("#div-charts", {name: "home"});
            $("#btn-add").click(this.showAddDlg.bind(this));
            $("#btn-add-chart").click(this.addChart.bind(this));
            $("#btn-save").click(() => {
                this.panel.save();
            });
        }

        private showAddDlg() {
            let $panel = $("#nav-charts");
            $panel.find("label.panel-block").remove();

            // load charts
            $ajax.get(`/system/chart/query`, {dashboard: "home"}).json((charts: any) => {
                for (let i = 0; i < charts.length; i++) {
                    let c = charts[i];
                    $panel.append(`<label class="panel-block">
          <input type="checkbox" value="${c.name}" data-index="${i}">${c.name}: ${c.title}
        </label>`);
                }
                this.charts = charts;
                this.$charts = $panel.find("label.panel-block");
            });

            let dlg = new Modal("#dlg-add-chart");
            dlg.show();
        }

        private filterCharts(text: string) {
            if (!text) {
                this.$charts.show();
                return;
            }

            this.$charts.each((i, elem) => {
                let $elem = $(elem);
                let texts: string[] = [
                    this.charts[i].name.toLowerCase(),
                    this.charts[i].title.toLowerCase(),
                    this.charts[i].desc.toLowerCase(),
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

        private addChart() {
            this.$charts.each((i, e) => {
                if ($(e).find(":checked").length > 0) {
                    let c = this.charts[i];
                    this.panel.addGraph(c);
                }
            });
            Modal.close();
        }
    }
}