///<reference path="../../core/core.ts" />
namespace Swirl.Service.Template {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class ListPage {
        constructor() {
            // bind events
            let dispatcher = Dispatcher.bind("#table-items");
            dispatcher.on("delete-template", this.deleteTemplate.bind(this));
        }

        private deleteTemplate(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let id = $tr.data("id");
            let name = $tr.find("td:first").text();
            Modal.confirm(`Are you sure to remove template: <strong>${name}</strong>?`, "Delete template", (dlg, e) => {
                $ajax.post("delete", { id: id }).trigger(e.target).encoder("form").json<AjaxResult>(() => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }
    }
}