import 'dart:convert';

import 'package:http/http.dart' as http;

/// Cliente mínimo da API Nest.
class KvoipApi {
  KvoipApi(this.baseUrl);

  final String baseUrl;

  Future<Map<String, dynamic>> login({
    required String email,
    required String password,
  }) async {
    final uri = Uri.parse('$baseUrl/auth/login');
    final res = await http.post(
      uri,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'email': email, 'password': password}),
    );
    if (res.statusCode < 200 || res.statusCode >= 300) {
      throw Exception(_message(res) ?? 'Falha no login (${res.statusCode})');
    }
    return jsonDecode(res.body) as Map<String, dynamic>;
  }

  String? _message(http.Response res) {
    try {
      final body = jsonDecode(res.body);
      if (body is Map && body['message'] != null) {
        return body['message'].toString();
      }
    } catch (_) {}
    return null;
  }
}
