///<reference path="core/core.ts" />
///<reference path="core/graph.ts" />
namespace Swirl {
    import Modal = Swirl.Core.Modal;
    import GraphPanel = Swirl.Core.GraphPanel;

    export class IndexPage {
        private panel: GraphPanel;

        constructor() {
            this.panel = new GraphPanel($("#div-charts").children("div"), {name: "home"});
            $("#btn-add").click(() => {
                Modal.alert("Coming soon...");
            });
        }
    }
}