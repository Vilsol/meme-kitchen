export interface RenderParams {
  context: CanvasRenderingContext2D;
  width: number;
  height: number;
}

export type RenderFunc = (params: RenderParams) => void