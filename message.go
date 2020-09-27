package ws

var _ Message = &BaseMessage{}

func NewBaseMessage(messageType MessageType, options ...MessageOption) *BaseMessage {
	var args = &MessageOptions{
		isSenderBroadcastSkipped: false,
	}

	var opt MessageOption
	for _, opt = range options {
		opt(args)
	}

	return &BaseMessage{
		MessageType:              messageType,
		isSenderBroadcastSkipped: args.isSenderBroadcastSkipped,
	}
}

func SkipSenderBroadcastSkipped() MessageOption {
	return func(options *MessageOptions) {
		options.isSenderBroadcastSkipped = true
	}
}

func (v *BaseMessage) GetType() MessageType {
	return v.MessageType
}

func (v *BaseMessage) IsSenderBroadcastSkipped() bool {
	return v.isSenderBroadcastSkipped
}
