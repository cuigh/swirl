///<reference path="../core/core.ts" />
namespace Swirl.Node {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class ListPage {
        constructor() {
            let dispatcher = Dispatcher.bind("#table-items");
            dispatcher.on("delete-node", this.deleteNode.bind(this));
        }

        private deleteNode(e: JQueryEventObject) {
            let $btn = $(e.target);
            let $tr = $btn.closest("tr");
            let id =$btn.val();
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to remove node: <strong>${name}</strong>?`, "Delete node", (dlg, e) => {
                $ajax.post("delete", {id: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })    
            });
        }
    }
}