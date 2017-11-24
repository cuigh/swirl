///<reference path="../core/core.ts" />
namespace Swirl.Perm {
    import Dispatcher = Swirl.Core.Dispatcher;
    import Modal = Swirl.Core.Modal;

    export class EditPage {
        constructor() {
            // bind events
            $("#txt-query").keydown(this.searchUser);
            $("#btn-add-user").click(this.addUser);
            Dispatcher.bind("#div-users").on("delete-user", this.deleteUser.bind(this));
        }

        private deleteUser(e: JQueryEventObject) {
            $(e.target).closest("div.control").remove();
        }

        private searchUser(e: JQueryEventObject) {
            if (e.keyCode == 13) {
                let query = $.trim($(e.target).val());
                if (query.length == 0) {
                    return;
                }
                $ajax.post("/system/user/search", {query: query}).encoder("form").json((users: any) => {
                    let $panel = $("#nav-users");
                    $panel.find("label.panel-block").remove();
                    for (let user of users) {
                        $panel.append(`<label class="panel-block">
          <input type="checkbox" value="${user.id}" data-name="${user.name}"> ${user.name}
        </label>`);
                    }
                });
            }
        }

        private addUser() {
            let users: { [index: string]: boolean } = {};
            $("#div-users").find("input").each((i, e) => {
                users[$(e).val()] = true;
            });

            let $panel = $("#nav-users");
            $panel.find("input:checked").each((i, e) => {
                let $el = $(e);
                if (users[$el.val()]) {
                    return;
                }

                $("#div-users").append(`<div class="control">
            <div class="tags has-addons">
              <span class="tag is-info">${$el.data("name")}</span>
              <a class="tag is-delete" data-action="delete-user"></a>
              <input name="users[]" value="${$el.val()}" type="hidden">
            </div>`);
            });

            Modal.close();
            $panel.find("label.panel-block").remove();
        }
    }
}