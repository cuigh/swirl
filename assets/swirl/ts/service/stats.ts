///<reference path="../core/core.ts" />
///<reference path="../core/chart.ts" />
namespace Swirl.Service {
    import Modal = Swirl.Core.Modal;
    import ChartDashboard = Swirl.Core.ChartDashboard;

    export class StatsPage {
        private dashboard: ChartDashboard;

        constructor() {
            let $cb_time = $("#cb-time");
            if ($cb_time.length == 0) {
                return;
            }

            this.dashboard = new ChartDashboard("#div-charts", window.charts, {
                name: "service",
                key: $("#h2-service-name").text()
            });

            // $("#btn-add").click(() => {
            //     Modal.alert("Coming soon...");
            // });
            // $("#btn-save").click(() => {
            //     Modal.alert("Coming soon...");
            // });
            $cb_time.change(e => {
                this.dashboard.setPeriod($(e.target).val());
            });
            $("#cb-refresh").change(e => {
                if ($(e.target).prop("checked")) {
                    this.dashboard.refresh();
                } else {
                    this.dashboard.stop();
                }
            });
        }
    }
}