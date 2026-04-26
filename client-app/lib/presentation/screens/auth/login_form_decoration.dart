import 'package:flutter/material.dart';

final class LoginFormDecoration {
  LoginFormDecoration._();

  static final outlineBorder = OutlineInputBorder(
    borderRadius: BorderRadius.circular(12),
  );

  static InputDecoration field({
    required String labelText,
    required String hintText,
    Widget? prefixIcon,
    Widget? suffixIcon,
  }) {
    return InputDecoration(
      labelText: labelText,
      hintText: hintText,
      prefixIcon: prefixIcon,
      suffixIcon: suffixIcon,
      border: outlineBorder,
    );
  }
}
