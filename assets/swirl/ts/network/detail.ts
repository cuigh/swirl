///<reference path="../core/core.ts" />
namespace Swirl.Network {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Dispatcher = Swirl.Core.Dispatcher;

    export class DetailPage {
        constructor() {
            let dispatcher = Dispatcher.bind("#table-containers");
            dispatcher.on("disconnect", this.disconnect.bind(this));
        }

        private disconnect(e: JQueryEventObject) {
            let $btn = $(e.target);
            let $tr = $btn.closest("tr");
            let id = $btn.val();
            let name = $tr.find("td:first").text().trim();
            Modal.confirm(`Are you sure to disconnect container: <strong>${name}</strong>?`, "Disconnect container", (dlg, e) => {
                $ajax.post("disconnect", {container: id}).trigger($btn).encoder("form").json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })    
            });
        }
    }
}