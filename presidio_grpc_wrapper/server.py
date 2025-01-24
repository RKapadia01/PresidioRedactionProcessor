from concurrent import futures
import logging

import grpc
import presidio_pb2
import presidio_pb2_grpc

from presidio_analyzer import AnalyzerEngine
from presidio_anonymizer import AnonymizerEngine
from presidio_anonymizer.entities import RecognizerResult
from presidio_anonymizer.entities import OperatorConfig

analyzer = AnalyzerEngine()
anonymizer = AnonymizerEngine()

class Server(presidio_pb2_grpc.PresidioRedactionProcessorServicer):
    def Analyze(self, request, context):
        recognizer_results = analyzer.analyze(
            text=request.text,
            language=request.language,
            entities=request.entities,
            score_threshold=request.score_threshold,
            context=request.context
        )

        responses = presidio_pb2.PresidioAnalyzerResponses()

        for r in recognizer_results:
            single_response = responses.analyzer_results.add()
            single_response.entity_type = r.entity_type
            single_response.score = r.score
            single_response.start = r.start
            single_response.end = r.end

        return responses

    def Anonymize(self, request, context):
        recognizer_results_py = []
        for proto_result in request.analyzer_results:
            py_result = RecognizerResult(
                entity_type=proto_result.entity_type,
                start=proto_result.start,
                end=proto_result.end,
                score=proto_result.score
            )
            recognizer_results_py.append(py_result)

        operators_py = {}
        for key, value in request.anonymizers.items():
            operator = OperatorConfig.from_json({
                "type": value.type,
                "new_value": value.new_value,
                "masking_char": value.masking_char,
                "chars_to_mask": value.chars_to_mask,
                "from_end": value.from_end,
                "hash_type": value.hash_type,
                "key": value.key
            })
            operators_py[key] = operator

        anonymizer_result = anonymizer.anonymize(
            text=request.text,
            operators=operators_py,
            analyzer_results=recognizer_results_py
        )
        
        return presidio_pb2.PresidioAnonymizerResponse(
            text=anonymizer_result.text,
        )

    def AnalyzeAndAnonymize(self, request, context):
        analyzer_results = analyzer.analyze(
            text=request.text, language=request.language,
            entities=request.entities, score_threshold=request.score_threshold,
            context=request.context)
        anonymizer_results = anonymizer.anonymize(
            text=request.text, analyzer_results=analyzer_results)
        return presidio_pb2.PresidioAnonymizerResponse(text=anonymizer_results)


def serve():
    port = "50051"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    presidio_pb2_grpc.add_PresidioRedactionProcessorServicer_to_server(Server(), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Server started, listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()