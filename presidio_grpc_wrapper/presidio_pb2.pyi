from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PresidioAnalyzerRequest(_message.Message):
    __slots__ = ("text", "language", "score_threshold", "entities", "context")
    TEXT_FIELD_NUMBER: _ClassVar[int]
    LANGUAGE_FIELD_NUMBER: _ClassVar[int]
    SCORE_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    ENTITIES_FIELD_NUMBER: _ClassVar[int]
    CONTEXT_FIELD_NUMBER: _ClassVar[int]
    text: str
    language: str
    score_threshold: float
    entities: _containers.RepeatedScalarFieldContainer[str]
    context: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, text: _Optional[str] = ..., language: _Optional[str] = ..., score_threshold: _Optional[float] = ..., entities: _Optional[_Iterable[str]] = ..., context: _Optional[_Iterable[str]] = ...) -> None: ...

class PresidioAnalyzerResponses(_message.Message):
    __slots__ = ("analyzer_results",)
    ANALYZER_RESULTS_FIELD_NUMBER: _ClassVar[int]
    analyzer_results: _containers.RepeatedCompositeFieldContainer[PresidioAnalyzerResponse]
    def __init__(self, analyzer_results: _Optional[_Iterable[_Union[PresidioAnalyzerResponse, _Mapping]]] = ...) -> None: ...

class PresidioAnalyzerResponse(_message.Message):
    __slots__ = ("start", "end", "score", "entity_type")
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    ENTITY_TYPE_FIELD_NUMBER: _ClassVar[int]
    start: int
    end: int
    score: float
    entity_type: str
    def __init__(self, start: _Optional[int] = ..., end: _Optional[int] = ..., score: _Optional[float] = ..., entity_type: _Optional[str] = ...) -> None: ...

class PresidioAnonymizerRequest(_message.Message):
    __slots__ = ("text", "anonymizers", "analyzer_results")
    class AnonymizersEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: PresidioAnonymizer
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[PresidioAnonymizer, _Mapping]] = ...) -> None: ...
    TEXT_FIELD_NUMBER: _ClassVar[int]
    ANONYMIZERS_FIELD_NUMBER: _ClassVar[int]
    ANALYZER_RESULTS_FIELD_NUMBER: _ClassVar[int]
    text: str
    anonymizers: _containers.MessageMap[str, PresidioAnonymizer]
    analyzer_results: _containers.RepeatedCompositeFieldContainer[PresidioAnalyzerResponse]
    def __init__(self, text: _Optional[str] = ..., anonymizers: _Optional[_Mapping[str, PresidioAnonymizer]] = ..., analyzer_results: _Optional[_Iterable[_Union[PresidioAnalyzerResponse, _Mapping]]] = ...) -> None: ...

class PresidioAnonymizer(_message.Message):
    __slots__ = ("type", "new_value", "masking_char", "chars_to_mask", "from_end", "hash_type", "key")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    NEW_VALUE_FIELD_NUMBER: _ClassVar[int]
    MASKING_CHAR_FIELD_NUMBER: _ClassVar[int]
    CHARS_TO_MASK_FIELD_NUMBER: _ClassVar[int]
    FROM_END_FIELD_NUMBER: _ClassVar[int]
    HASH_TYPE_FIELD_NUMBER: _ClassVar[int]
    KEY_FIELD_NUMBER: _ClassVar[int]
    type: str
    new_value: str
    masking_char: str
    chars_to_mask: int
    from_end: bool
    hash_type: str
    key: str
    def __init__(self, type: _Optional[str] = ..., new_value: _Optional[str] = ..., masking_char: _Optional[str] = ..., chars_to_mask: _Optional[int] = ..., from_end: bool = ..., hash_type: _Optional[str] = ..., key: _Optional[str] = ...) -> None: ...

class PresidioAnonymizerResponse(_message.Message):
    __slots__ = ("operation", "entity_type", "start", "end", "text")
    OPERATION_FIELD_NUMBER: _ClassVar[int]
    ENTITY_TYPE_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    operation: str
    entity_type: str
    start: int
    end: int
    text: str
    def __init__(self, operation: _Optional[str] = ..., entity_type: _Optional[str] = ..., start: _Optional[int] = ..., end: _Optional[int] = ..., text: _Optional[str] = ...) -> None: ...

class PresidioAnalyzerAnomymizerRequest(_message.Message):
    __slots__ = ("text", "language", "score_threshold", "entities", "context", "anonymizers")
    class AnonymizersEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: PresidioAnonymizer
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[PresidioAnonymizer, _Mapping]] = ...) -> None: ...
    TEXT_FIELD_NUMBER: _ClassVar[int]
    LANGUAGE_FIELD_NUMBER: _ClassVar[int]
    SCORE_THRESHOLD_FIELD_NUMBER: _ClassVar[int]
    ENTITIES_FIELD_NUMBER: _ClassVar[int]
    CONTEXT_FIELD_NUMBER: _ClassVar[int]
    ANONYMIZERS_FIELD_NUMBER: _ClassVar[int]
    text: str
    language: str
    score_threshold: float
    entities: _containers.RepeatedScalarFieldContainer[str]
    context: _containers.RepeatedScalarFieldContainer[str]
    anonymizers: _containers.MessageMap[str, PresidioAnonymizer]
    def __init__(self, text: _Optional[str] = ..., language: _Optional[str] = ..., score_threshold: _Optional[float] = ..., entities: _Optional[_Iterable[str]] = ..., context: _Optional[_Iterable[str]] = ..., anonymizers: _Optional[_Mapping[str, PresidioAnonymizer]] = ...) -> None: ...
