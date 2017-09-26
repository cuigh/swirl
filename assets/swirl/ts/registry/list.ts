///<reference path="../core/core.ts" />
namespace Swirl.Registry {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class ListPage {
        constructor() {
            let dispatcher = Dispatcher.bind("#table-items");
            dispatcher.on("delete-registry", this.deleteRegistry.bind(this));
            dispatcher.on("edit-registry", this.editRegistry.bind(this));
        }

        private deleteRegistry(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let id = $tr.data("id");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to remove registry: <strong>${name}</strong>?`, "Delete registry", (dlg, e) => {
                $ajax.post("delete", {id: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {                
                    $tr.remove();
                    dlg.close();
                })    
            });
        }

        private editRegistry(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let dlg = new Modal("#dlg-edit");
            dlg.find("input[name=id]").val($tr.data("id"));
            dlg.find("input[name=name]").val($tr.find("td:first").text().trim());
            dlg.find("input[name=url]").val($tr.find("td:eq(1)").text().trim());
            dlg.find("input[name=username]").val($tr.find("td:eq(2)").text().trim());
            dlg.show();
        }
    }
}