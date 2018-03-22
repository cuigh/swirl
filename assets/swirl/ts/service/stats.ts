///<reference path="../core/core.ts" />
///<reference path="../core/graph.ts" />
namespace Swirl.Service {
    import Modal = Swirl.Core.Modal;
    import GraphPanel = Swirl.Core.GraphPanel;

    export class StatsPage {
        private panel: GraphPanel;

        constructor() {
            let $cb_time = $("#cb-time");
            if ($cb_time.length == 0) {
                return;
            }

            this.panel = new GraphPanel($("#div-charts").children("div"), {
                name: "service",
                id: $("#h2-service-name").text()
            });

            $("#btn-add").click(() => {
                Modal.alert("Coming soon...");
            });
            $cb_time.change(e => {
                this.panel.setTime($(e.target).val());
            });
            $("#cb-refresh").change(e => {
                if ($(e.target).prop("checked")) {
                    this.panel.refresh();
                } else {
                    this.panel.stop();
                }
            });
        }
    }
}