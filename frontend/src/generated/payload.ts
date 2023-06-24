/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal.js";

export enum HorizontalAlign {
  CENTER = 0,
  LEFT = 1,
  RIGHT = 2,
  UNRECOGNIZED = -1,
}

export function horizontalAlignFromJSON(object: any): HorizontalAlign {
  switch (object) {
    case 0:
    case "CENTER":
      return HorizontalAlign.CENTER;
    case 1:
    case "LEFT":
      return HorizontalAlign.LEFT;
    case 2:
    case "RIGHT":
      return HorizontalAlign.RIGHT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return HorizontalAlign.UNRECOGNIZED;
  }
}

export function horizontalAlignToJSON(object: HorizontalAlign): string {
  switch (object) {
    case HorizontalAlign.CENTER:
      return "CENTER";
    case HorizontalAlign.LEFT:
      return "LEFT";
    case HorizontalAlign.RIGHT:
      return "RIGHT";
    case HorizontalAlign.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum VerticalAlign {
  MIDDLE = 0,
  TOP = 1,
  BOTTOM = 2,
  UNRECOGNIZED = -1,
}

export function verticalAlignFromJSON(object: any): VerticalAlign {
  switch (object) {
    case 0:
    case "MIDDLE":
      return VerticalAlign.MIDDLE;
    case 1:
    case "TOP":
      return VerticalAlign.TOP;
    case 2:
    case "BOTTOM":
      return VerticalAlign.BOTTOM;
    case -1:
    case "UNRECOGNIZED":
    default:
      return VerticalAlign.UNRECOGNIZED;
  }
}

export function verticalAlignToJSON(object: VerticalAlign): string {
  switch (object) {
    case VerticalAlign.MIDDLE:
      return "MIDDLE";
    case VerticalAlign.TOP:
      return "TOP";
    case VerticalAlign.BOTTOM:
      return "BOTTOM";
    case VerticalAlign.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Payload {
  version: number;
  template: number;
  text: Text[];
}

export interface Text {
  template_text?: number | undefined;
  text?: string | undefined;
  x?: number | undefined;
  y?: number | undefined;
  width?: number | undefined;
  height?: number | undefined;
  font?: string | undefined;
  size?: number | undefined;
  unfilled?: boolean | undefined;
  fill_color?: string | undefined;
  stroke_color?: string | undefined;
  stroke?: number | undefined;
  horizontal_align?: HorizontalAlign | undefined;
  vertical_align?: VerticalAlign | undefined;
}

export interface TemplateText {
  text?: string | undefined;
  x: number;
  y: number;
  width: number;
  height: number;
  font: string;
  size: number;
  unfilled: boolean;
  fill_color: string;
  stroke_color: string;
  stroke: number;
  horizontal_align: HorizontalAlign;
  vertical_align: VerticalAlign;
}

function createBasePayload(): Payload {
  return { version: 0, template: 0, text: [] };
}

export const Payload = {
  encode(message: Payload, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.version !== 0) {
      writer.uint32(8).uint32(message.version);
    }
    if (message.template !== 0) {
      writer.uint32(16).uint64(message.template);
    }
    for (const v of message.text) {
      Text.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Payload {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePayload();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.version = reader.uint32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.template = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.text.push(Text.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Payload {
    return {
      version: isSet(object.version) ? Number(object.version) : 0,
      template: isSet(object.template) ? Number(object.template) : 0,
      text: Array.isArray(object?.text) ? object.text.map((e: any) => Text.fromJSON(e)) : [],
    };
  },

  toJSON(message: Payload): unknown {
    const obj: any = {};
    message.version !== undefined && (obj.version = Math.round(message.version));
    message.template !== undefined && (obj.template = Math.round(message.template));
    if (message.text) {
      obj.text = message.text.map((e) => e ? Text.toJSON(e) : undefined);
    } else {
      obj.text = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Payload>, I>>(base?: I): Payload {
    return Payload.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Payload>, I>>(object: I): Payload {
    const message = createBasePayload();
    message.version = object.version ?? 0;
    message.template = object.template ?? 0;
    message.text = object.text?.map((e) => Text.fromPartial(e)) || [];
    return message;
  },
};

function createBaseText(): Text {
  return {
    template_text: undefined,
    text: undefined,
    x: undefined,
    y: undefined,
    width: undefined,
    height: undefined,
    font: undefined,
    size: undefined,
    unfilled: undefined,
    fill_color: undefined,
    stroke_color: undefined,
    stroke: undefined,
    horizontal_align: undefined,
    vertical_align: undefined,
  };
}

export const Text = {
  encode(message: Text, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.template_text !== undefined) {
      writer.uint32(8).uint32(message.template_text);
    }
    if (message.text !== undefined) {
      writer.uint32(18).string(message.text);
    }
    if (message.x !== undefined) {
      writer.uint32(24).uint32(message.x);
    }
    if (message.y !== undefined) {
      writer.uint32(32).uint32(message.y);
    }
    if (message.width !== undefined) {
      writer.uint32(40).uint32(message.width);
    }
    if (message.height !== undefined) {
      writer.uint32(48).uint32(message.height);
    }
    if (message.font !== undefined) {
      writer.uint32(58).string(message.font);
    }
    if (message.size !== undefined) {
      writer.uint32(64).uint32(message.size);
    }
    if (message.unfilled !== undefined) {
      writer.uint32(72).bool(message.unfilled);
    }
    if (message.fill_color !== undefined) {
      writer.uint32(82).string(message.fill_color);
    }
    if (message.stroke_color !== undefined) {
      writer.uint32(90).string(message.stroke_color);
    }
    if (message.stroke !== undefined) {
      writer.uint32(96).uint32(message.stroke);
    }
    if (message.horizontal_align !== undefined) {
      writer.uint32(104).int32(message.horizontal_align);
    }
    if (message.vertical_align !== undefined) {
      writer.uint32(112).int32(message.vertical_align);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Text {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.template_text = reader.uint32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.text = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.x = reader.uint32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.y = reader.uint32();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.width = reader.uint32();
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.height = reader.uint32();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.font = reader.string();
          continue;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.size = reader.uint32();
          continue;
        case 9:
          if (tag !== 72) {
            break;
          }

          message.unfilled = reader.bool();
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.fill_color = reader.string();
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.stroke_color = reader.string();
          continue;
        case 12:
          if (tag !== 96) {
            break;
          }

          message.stroke = reader.uint32();
          continue;
        case 13:
          if (tag !== 104) {
            break;
          }

          message.horizontal_align = reader.int32() as any;
          continue;
        case 14:
          if (tag !== 112) {
            break;
          }

          message.vertical_align = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Text {
    return {
      template_text: isSet(object.template_text) ? Number(object.template_text) : undefined,
      text: isSet(object.text) ? String(object.text) : undefined,
      x: isSet(object.x) ? Number(object.x) : undefined,
      y: isSet(object.y) ? Number(object.y) : undefined,
      width: isSet(object.width) ? Number(object.width) : undefined,
      height: isSet(object.height) ? Number(object.height) : undefined,
      font: isSet(object.font) ? String(object.font) : undefined,
      size: isSet(object.size) ? Number(object.size) : undefined,
      unfilled: isSet(object.unfilled) ? Boolean(object.unfilled) : undefined,
      fill_color: isSet(object.fill_color) ? String(object.fill_color) : undefined,
      stroke_color: isSet(object.stroke_color) ? String(object.stroke_color) : undefined,
      stroke: isSet(object.stroke) ? Number(object.stroke) : undefined,
      horizontal_align: isSet(object.horizontal_align) ? horizontalAlignFromJSON(object.horizontal_align) : undefined,
      vertical_align: isSet(object.vertical_align) ? verticalAlignFromJSON(object.vertical_align) : undefined,
    };
  },

  toJSON(message: Text): unknown {
    const obj: any = {};
    message.template_text !== undefined && (obj.template_text = Math.round(message.template_text));
    message.text !== undefined && (obj.text = message.text);
    message.x !== undefined && (obj.x = Math.round(message.x));
    message.y !== undefined && (obj.y = Math.round(message.y));
    message.width !== undefined && (obj.width = Math.round(message.width));
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.font !== undefined && (obj.font = message.font);
    message.size !== undefined && (obj.size = Math.round(message.size));
    message.unfilled !== undefined && (obj.unfilled = message.unfilled);
    message.fill_color !== undefined && (obj.fill_color = message.fill_color);
    message.stroke_color !== undefined && (obj.stroke_color = message.stroke_color);
    message.stroke !== undefined && (obj.stroke = Math.round(message.stroke));
    message.horizontal_align !== undefined && (obj.horizontal_align = message.horizontal_align !== undefined
      ? horizontalAlignToJSON(message.horizontal_align)
      : undefined);
    message.vertical_align !== undefined && (obj.vertical_align = message.vertical_align !== undefined
      ? verticalAlignToJSON(message.vertical_align)
      : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<Text>, I>>(base?: I): Text {
    return Text.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<Text>, I>>(object: I): Text {
    const message = createBaseText();
    message.template_text = object.template_text ?? undefined;
    message.text = object.text ?? undefined;
    message.x = object.x ?? undefined;
    message.y = object.y ?? undefined;
    message.width = object.width ?? undefined;
    message.height = object.height ?? undefined;
    message.font = object.font ?? undefined;
    message.size = object.size ?? undefined;
    message.unfilled = object.unfilled ?? undefined;
    message.fill_color = object.fill_color ?? undefined;
    message.stroke_color = object.stroke_color ?? undefined;
    message.stroke = object.stroke ?? undefined;
    message.horizontal_align = object.horizontal_align ?? undefined;
    message.vertical_align = object.vertical_align ?? undefined;
    return message;
  },
};

function createBaseTemplateText(): TemplateText {
  return {
    text: undefined,
    x: 0,
    y: 0,
    width: 0,
    height: 0,
    font: "",
    size: 0,
    unfilled: false,
    fill_color: "",
    stroke_color: "",
    stroke: 0,
    horizontal_align: 0,
    vertical_align: 0,
  };
}

export const TemplateText = {
  encode(message: TemplateText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.text !== undefined) {
      writer.uint32(10).string(message.text);
    }
    if (message.x !== 0) {
      writer.uint32(16).uint32(message.x);
    }
    if (message.y !== 0) {
      writer.uint32(24).uint32(message.y);
    }
    if (message.width !== 0) {
      writer.uint32(32).uint32(message.width);
    }
    if (message.height !== 0) {
      writer.uint32(40).uint32(message.height);
    }
    if (message.font !== "") {
      writer.uint32(50).string(message.font);
    }
    if (message.size !== 0) {
      writer.uint32(56).uint32(message.size);
    }
    if (message.unfilled === true) {
      writer.uint32(64).bool(message.unfilled);
    }
    if (message.fill_color !== "") {
      writer.uint32(74).string(message.fill_color);
    }
    if (message.stroke_color !== "") {
      writer.uint32(82).string(message.stroke_color);
    }
    if (message.stroke !== 0) {
      writer.uint32(88).uint32(message.stroke);
    }
    if (message.horizontal_align !== 0) {
      writer.uint32(96).int32(message.horizontal_align);
    }
    if (message.vertical_align !== 0) {
      writer.uint32(104).int32(message.vertical_align);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TemplateText {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTemplateText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.text = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.x = reader.uint32();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.y = reader.uint32();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.width = reader.uint32();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.height = reader.uint32();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.font = reader.string();
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.size = reader.uint32();
          continue;
        case 8:
          if (tag !== 64) {
            break;
          }

          message.unfilled = reader.bool();
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.fill_color = reader.string();
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.stroke_color = reader.string();
          continue;
        case 11:
          if (tag !== 88) {
            break;
          }

          message.stroke = reader.uint32();
          continue;
        case 12:
          if (tag !== 96) {
            break;
          }

          message.horizontal_align = reader.int32() as any;
          continue;
        case 13:
          if (tag !== 104) {
            break;
          }

          message.vertical_align = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TemplateText {
    return {
      text: isSet(object.text) ? String(object.text) : undefined,
      x: isSet(object.x) ? Number(object.x) : 0,
      y: isSet(object.y) ? Number(object.y) : 0,
      width: isSet(object.width) ? Number(object.width) : 0,
      height: isSet(object.height) ? Number(object.height) : 0,
      font: isSet(object.font) ? String(object.font) : "",
      size: isSet(object.size) ? Number(object.size) : 0,
      unfilled: isSet(object.unfilled) ? Boolean(object.unfilled) : false,
      fill_color: isSet(object.fill_color) ? String(object.fill_color) : "",
      stroke_color: isSet(object.stroke_color) ? String(object.stroke_color) : "",
      stroke: isSet(object.stroke) ? Number(object.stroke) : 0,
      horizontal_align: isSet(object.horizontal_align) ? horizontalAlignFromJSON(object.horizontal_align) : 0,
      vertical_align: isSet(object.vertical_align) ? verticalAlignFromJSON(object.vertical_align) : 0,
    };
  },

  toJSON(message: TemplateText): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    message.x !== undefined && (obj.x = Math.round(message.x));
    message.y !== undefined && (obj.y = Math.round(message.y));
    message.width !== undefined && (obj.width = Math.round(message.width));
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.font !== undefined && (obj.font = message.font);
    message.size !== undefined && (obj.size = Math.round(message.size));
    message.unfilled !== undefined && (obj.unfilled = message.unfilled);
    message.fill_color !== undefined && (obj.fill_color = message.fill_color);
    message.stroke_color !== undefined && (obj.stroke_color = message.stroke_color);
    message.stroke !== undefined && (obj.stroke = Math.round(message.stroke));
    message.horizontal_align !== undefined && (obj.horizontal_align = horizontalAlignToJSON(message.horizontal_align));
    message.vertical_align !== undefined && (obj.vertical_align = verticalAlignToJSON(message.vertical_align));
    return obj;
  },

  create<I extends Exact<DeepPartial<TemplateText>, I>>(base?: I): TemplateText {
    return TemplateText.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<TemplateText>, I>>(object: I): TemplateText {
    const message = createBaseTemplateText();
    message.text = object.text ?? undefined;
    message.x = object.x ?? 0;
    message.y = object.y ?? 0;
    message.width = object.width ?? 0;
    message.height = object.height ?? 0;
    message.font = object.font ?? "";
    message.size = object.size ?? 0;
    message.unfilled = object.unfilled ?? false;
    message.fill_color = object.fill_color ?? "";
    message.stroke_color = object.stroke_color ?? "";
    message.stroke = object.stroke ?? 0;
    message.horizontal_align = object.horizontal_align ?? 0;
    message.vertical_align = object.vertical_align ?? 0;
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
