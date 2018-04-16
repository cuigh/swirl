///<reference path="../core/core.ts" />
namespace Swirl.Stack {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class ListPage {
        constructor() {
            let dispatcher = Dispatcher.bind("#table-items");
            dispatcher.on("deploy-stack", this.deployStack.bind(this));
            dispatcher.on("shutdown-stack", this.shutdownStack.bind(this));
            dispatcher.on("delete-stack", this.deleteStack.bind(this));
        }

        private deployStack(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to deploy stack: <strong>${name}</strong>?`, "Deploy stack", (dlg, e) => {
                $ajax.post(`${name}/deploy`).trigger(e.target).json<AjaxResult>(r => {
                    location.reload();
                })
            });
        }

        private shutdownStack(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to shutdown stack: <strong>${name}</strong>?`, "Shutdown stack", (dlg, e) => {
                $ajax.post(`${name}/shutdown`).trigger(e.target).json<AjaxResult>(r => {
                    location.reload();
                })
            });
        }

        private deleteStack(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to remove archive: <strong>${name}</strong>?`, "Delete stack", (dlg, e) => {
                $ajax.post(`${name}/delete`).trigger(e.target).json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }
    }
}