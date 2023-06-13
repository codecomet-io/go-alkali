package trace

import (
	"strings"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
)

/*
import (
	"encoding/json"
	"github.com/moby/buildkit/client"
)

type V = client.VertexStatus

func Ingest(w []byte) *client.VertexStatus {
	solved := &client.SolveStatus{}
	_ = json.Unmarshal(w, solved)
	// o.op.Run.Trace.Bytes()

}
*/

func Ingest(state llb.State, statusSlice []*client.SolveStatus) {
	//marshmallow := state.Marshal(context.TODO())

	for _, status := range statusSlice {
		if status.Vertexes != nil {
			for _, vtx := range status.Vertexes {
				vertex(vtx)
			}
		}
	}
}

func vertex(vtx *client.Vertex) {
	// Some actions are hidden away - either CodeComet internal shenanigans, or actions authors who want to hide their own internal dance
	if vtx.ProgressGroup != nil && vtx.ProgressGroup.Weak {
		return
	}

	// Currently, BK leaks internal operations. The right solution is to finish replacing the default client with our own. Short term, very dirty hack by ignoring anything that starts with "[auth] "
	if strings.HasPrefix(vtx.Name, "[auth] ") {
		return
	}
	/*

		if (!this.actionsObject[vertice.Digest]) {
			const action = <model.BuildAction>{
			id     : createId('html_id'),
				name   : vertice.Name,
					digest : vertice.Digest,
					cached : false,
					status : ActionStatus.Ignored,
			};

			if (vertice.Inputs) {
				action.buildParents = vertice.Inputs;
			}

			this.actionsObject[vertice.Digest] = action;
		}

		if (vertice.Started){
			this.actionsObject[vertice.Digest].started = Date.parse(vertice.Started);

			if (!this.started || this.actionsObject[vertice.Digest].started < this.started) {
				this.started = this.actionsObject[vertice.Digest].started;
			}

			this.actionsObject[vertice.Digest].status = ActionStatus.Started;
		}

		if (vertice.Completed) {
			this.actionsObject[vertice.Digest].completed = Date.parse(vertice.Completed);
			this.actionsObject[vertice.Digest].runtime = this.actionsObject[vertice.Digest].completed - this.actionsObject[vertice.Digest].started;
			this.actionsObject[vertice.Digest].status = ActionStatus.Completed;
		}

		if (vertice.Error) {
			this.actionsObject[vertice.Digest].error = vertice.Error;
			this.actionsObject[vertice.Digest].status = ActionStatus.Errored;

			if (vertice.Error.match(/did not complete successfully: exit code: 137: context canceled: context canceled$/)) {
			this.actionsObject[vertice.Digest].status = ActionStatus.Cancelled;
			}
		}

		if (vertice.Cached) {
			this.actionsObject[vertice.Digest].cached = true;
			this.actionsObject[vertice.Digest].status = ActionStatus.Cached;
		}


	*/
}
