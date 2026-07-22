import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class KvoipTheme {
  static const brand = Color(0xFF0080FF);
  static const brandStrong = Color(0xFF0066D6);
  static const ink = Color(0xFF0B1B2B);
  static const muted = Color(0xFF5B6B7C);
  static const surface = Color(0xFFF7FBFF);

  static ThemeData light() {
    final base = ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: brand,
        primary: brand,
        surface: surface,
      ),
    );
    return base.copyWith(
      textTheme: GoogleFonts.sourceSans3TextTheme(base.textTheme).apply(
        bodyColor: ink,
        displayColor: ink,
      ),
      appBarTheme: AppBarTheme(
        backgroundColor: Colors.white,
        foregroundColor: ink,
        elevation: 0,
        titleTextStyle: GoogleFonts.montserrat(
          fontWeight: FontWeight.w700,
          fontSize: 18,
          color: ink,
        ),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: brand,
          foregroundColor: Colors.white,
          padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 18),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: Colors.white,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: const BorderSide(color: Color(0xFFD7E3F0)),
        ),
      ),
    );
  }
}
