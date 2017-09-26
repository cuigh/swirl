///<reference path="../core/core.ts" />
namespace Swirl.User {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class ListPage {
        constructor() {
            // bind events
            Dispatcher.bind("#table-items")
                .on("delete-user", this.deleteUser.bind(this))
                .on("block-user", this.blockUser.bind(this))
                .on("unblock-user", this.unblockUser.bind(this));
        }

        private deleteUser(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let id = $tr.data("id");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to remove user: <strong>${name}</strong>?`, "Delete user", (dlg, e) => {
                $ajax.post("delete", {id: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }

        private blockUser(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let id = $tr.data("id");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to block user: <strong>${name}</strong>?`, "Block user", (dlg, e) => {
                $ajax.post("block", {id: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    location.reload();
                })
            });
        }

        private unblockUser(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let id = $tr.data("id");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to unblock user: <strong>${name}</strong>?`, "Unblock user", (dlg, e) => {
                $ajax.post("unblock", {id: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    location.reload();
                })
            });
        }
    }
}