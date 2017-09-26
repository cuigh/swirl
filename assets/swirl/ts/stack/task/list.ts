///<reference path="../../core/core.ts" />
namespace Swirl.Stack.Task {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class ListPage {
        constructor() {
            let dispatcher = Dispatcher.bind("#table-items");
            dispatcher.on("delete-stack", this.deleteStack.bind(this));
        }

        private deleteStack(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to remove stack: <strong>${name}</strong>?`, "Delete stack", (dlg, e) => {
                $ajax.post("delete", {name: name}).trigger(e.target).encoder("form").json<AjaxResult>(r => {                
                    $tr.remove();
                    dlg.close();
                })    
            });
        }
    }
}