///<reference path="../../core/core.ts" />
namespace Swirl.Stack.Archive {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class ListPage {
        constructor() {
            let dispatcher = Dispatcher.bind("#table-items");
            dispatcher.on("deploy-archive", this.deployArchive.bind(this));
            dispatcher.on("delete-archive", this.deleteArchive.bind(this));
        }

        private deployArchive(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let id = $tr.data("id");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to deploy archive: <strong>${name}</strong>?`, "Deploy archive", (dlg, e) => {
                $ajax.post("deploy", {id: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    dlg.close();
                })
            });
        }

        private deleteArchive(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let id = $tr.data("id");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to remove archive: <strong>${name}</strong>?`, "Delete archive", (dlg, e) => {
                $ajax.post("delete", {id: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })    
            });
        }
    }
}