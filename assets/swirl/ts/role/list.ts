///<reference path="../core/core.ts" />
namespace Swirl.Role {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class ListPage {
        constructor() {
            // bind events
            Dispatcher.bind("#table-items")
                .on("delete-role", this.deleteUser.bind(this))
        }

        private deleteUser(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let id = $tr.data("id");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to remove role: <strong>${name}</strong>?`, "Delete role", (dlg, e) => {
                $ajax.post("delete", {id: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }
    }
}